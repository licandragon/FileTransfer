<template>
    <div class="relative z-10 flex justify-center">
        <div class="bg-[var(--bg-main)]
            border border-[rgba(200,169,110,0.15)]
            rounded-[14px] overflow-hidden
            w-full max-w-85 sm:max-w-95
            shadow-[0_0_0_1px_rgba(255,255,255,0.03),0_24px_48px_rgba(0,0,0,0.4)]
            animate-floatCard">

            <div class="flex gap-1.5 px-4 py-[0.7rem]
            border-b border-[rgba(255,255,255,0.05)]
            bg-[rgba(255,255,255,0.02)]">
                <span class="w-2.5 h-2.5 rounded-full bg-[rgba(255,255,255,0.12)]"></span>
                <span class="w-2.5 h-2.5 rounded-full bg-[rgba(255,255,255,0.12)]"></span>
                <span class="w-2.5 h-2.5 rounded-full bg-[rgba(255,255,255,0.12)]"></span>
            </div>

            <pre class="px-4 sm:px-6 py-5 text-[0.72rem] sm:text-[0.78rem]
            font-mono leading-[1.75]
            overflow-x-auto w-full">
                <code>
                    <span class="text-[#4a4a5a]">// Navigation Guard</span>
                    <span class="text-[#6eb4e8]">router</span>.<span class="text-[var(--primary)]">beforeEach</span>(<span class="text-[#b08ce8]">(to, from)</span> => {
                    <span class="text-[#6eb4e8]">const</span> { isAuthenticated } = <span class="text-[var(--primary)]">useAuth</span>()
                    <span class="text-[#6eb4e8]">if</span> (to.meta.<span class="text-[#e8c46e]">requiresAuth</span>
                    && !isAuthenticated.<span class="text-[#e8c46e]">value</span>) {
                        <span class="text-[#6eb4e8]">return</span> { <span class="text-[#e8c46e]">name</span>: <span class="text-[#6fcf97]">'login'</span> }
                    }
                    })
                </code>
            </pre>
        </div>
    </div>
</template>