import { ref } from "vue";
import axios from "axios";

export function useTransferApi() {
    const uploadProgress = ref(0);
    const isUploading = ref(false);
    const error = ref(null);

    // Instancia centralizada de Axios
    const api = axios.create({
        baseURL: "http://localhost:3000/api",
    });

    // 1. Iniciar transferencia
    async function createTransfer(metadata) {
        const { data } = await api.post("/transfer/init", {
            sender_email: metadata.senderEmail,
            subject_email: metadata.subject,
            message_email: metadata.message,
            recipients: metadata.recipients,
            total_files: metadata.totalFiles,
            expires_at: metadata.expiredAt,
            mode: metadata.mode
        });
        console.log("transferencia creada")
        return data.upload_token;
    }

    // 3. Consulta el estatus del la transferencia actual
    async function getTransferStatus(uploadToken) {
        try {
            const { data } = await api.get(`/transfer/${uploadToken}/status`);
            return data.completedIndices || []; // Retorna ej: [0, 1]
        } catch (err) {
            console.warn("No se pudo recuperar el estado de reintento, se asumirá limpio:", err);
            return [];
        }
    }
    // 3. Subir archivo con callback para calcular el progreso global
    async function uploadFile(uploadToken, file, fileIndex, onProgressCallback) {
        const formData = new FormData();
        formData.append("file", file);
        formData.append("file_index", fileIndex.toString());
        console.log("subiendo archivo:", formData)

        return await api.post(`/transfer/${uploadToken}/file`, formData, {
            onUploadProgress: (progressEvent) => {
                if (onProgressCallback) {
                    onProgressCallback(progressEvent.loaded);
                }
            },
        });
    }

    // 3. Finalizar subida
    async function completeTransfer(uploadToken) {
        const { data } = await api.patch(`/transfer/${uploadToken}/complete`);
        return data.download_token;
    }

    // 4. Obtener info para la vista de descarga
    async function getDownloadInfo(downloadToken) {
        const { data } = await api.get(`/download/${downloadToken}`);
        return data;
    }

    // 5. Descargar archivo específico
    async function downloadFile(downloadToken, fileIndex, fileName) {
        const response = await api({
            url: `/download/${downloadToken}/files/${fileIndex}`,
            method: "GET",
            responseType: "blob",
        });

        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", fileName);
        document.body.appendChild(link);
        link.click();
        link.remove();
        window.URL.revokeObjectURL(url); // Libera memoria del navegador
    }

    // Subida de archivos
    async function startTransfer(metadata, files, existingToken = null) {
        console.log("Start transfer archivos", files)
        isUploading.value = true;
        error.value = null;
        let uploadToken = existingToken;
        let completedIndices = [];

        try {
            //
            if (!uploadToken) {
                //Se genera transfer nuevo
                uploadProgress.value = 0;

                uploadToken = await createTransfer(metadata);
            } else {
                //En caso de reintento por error en la red o subida de archivo
                completedIndices = await getTransferStatus(uploadToken);
                console.log("Índices recuperados del servidor para reintento:", completedIndices);
            }

            // Calcular el peso total de todos los archivos combinados
            const totalBytesGlobal = files.reduce((acc, file) => acc + file.size, 0);
            const bytesSubidosPorArchivo = new Array(files.length).fill(0);

            for (let i = 0; i < files.length; i++) {
                if (completedIndices.includes(i)) {
                    console.log(`[Skip] Archivo index ${i} (${files[i].name}) ya existe.`);

                    bytesSubidosPorArchivo[i] = files[i].size; // 👈 Seteamos el peso real aquí

                    const totalBytesSubidosProcesados = bytesSubidosPorArchivo.reduce((acc, b) => acc + b, 0);
                    uploadProgress.value = Math.round((totalBytesSubidosProcesados * 100) / totalBytesGlobal);
                    continue;
                }

                    await uploadFile(
                        uploadToken,
                        files[i],
                        i,
                        (bytesSubidosEsteArchivo) => {
                            // Guardamos cuántos bytes lleva este archivo actual
                            bytesSubidosPorArchivo[i] = bytesSubidosEsteArchivo;

                            // Sumamos los bytes de todos los archivos para sacar el porcentaje real global
                            const totalBytesSubidosProcesados = bytesSubidosPorArchivo.reduce(
                                (acc, bytes) => acc + bytes,
                                0,
                            );
                            uploadProgress.value = Math.round(
                                (totalBytesSubidosProcesados * 100) / totalBytesGlobal,
                            );
                        },
                    );
                
            }
                const downloadToken = await completeTransfer(uploadToken);
                isUploading.value = false;
                return downloadToken;
            } catch (err) {
                error.value = err.response?.data?.error || err.message;

                throw err;
            }
        }

    return {
            uploadProgress,
            isUploading,
            error,
            startTransfer,
            getDownloadInfo,
            downloadFile,
        };
    }
