<template>
    <section class="w-[45%] bg-[#191819] text-gray-300 p-6 rounded-xl flex flex-col">
        <!-- Filtros -->
        <div class="flex flex-col gap-4 mb-6">
            <!-- Barra de búsqueda -->
            <div class="relative">
                <span class="absolute inset-y-0 left-0 flex items-center pl-3">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500" fill="none"
                        viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                </span>
                <input type="text" 
                    v-model="searchQuery" 
                    @input="debounceSearch" 
                    placeholder="Buscar emails..."
                    class="w-full pl-10 pr-10 py-3 rounded-full bg-[#202021] placeholder-gray-500 text-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                <!-- X -->
                <button v-if="hasAnyFilter" 
                    @click="clearFilters"
                    class="absolute inset-y-0 right-0 flex items-center pr-3">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500 hover:text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>
            <!-- Remitente y Destinatario -->
            <div v-if="searchQuery.trim()" class="flex gap-4">
                <input type="text" 
                    v-model="fromFilter" 
                    @input="debounceSearch" 
                    placeholder="Remitente"
                    class="w-full pl-4 pr-4 py-2 rounded-full bg-[#202021] placeholder-gray-500 text-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                <input type="text" 
                    v-model="toFilter" 
                    @input="debounceSearch" 
                    placeholder="Destinatario"
                    class="w-full pl-4 pr-4 py-2 rounded-full bg-[#202021] placeholder-gray-500 text-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500" />
            </div>
        </div>

        <div v-if="isLoading" class="flex justify-center items-center py-4">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>

        <div v-if="!isLoading && hasSearched && totalEmails === 0 && (searchQuery.trim() || fromFilter.trim() || toFilter.trim())" 
            class="flex justify-center items-center py-8 text-gray-400">
            <span>No se encontraron coincidencias</span>
        </div>

        <div class="space-y-4 flex-1">
            <div v-for="email in paginatedEmails" :key="email.id"
                class="flex items-center p-4 bg-[#202021] rounded-lg shadow-md cursor-pointer hover:bg-[#202021] transition-colors h-32"
                @click="selectEmail(email)">
                <div
                    class="flex-shrink-0 w-10 h-10 rounded-full bg-gray-600 flex items-center justify-center text-white font-bold mr-4">
                    {{ email.sender.charAt(0).toUpperCase() }}
                </div>
                <!-- Información del email -->
                <div class="flex-1">
                    <div class="flex justify-between items-center">
                        <span class="font-bold truncate" v-html="email.sender_preview"></span>
                        <span class="text-sm text-gray-400">{{ email.date }}</span>
                    </div>
                    <h3 class="text-m font-semibold truncate" v-html="email.subject_preview"></h3>
                    <p class="text-sm text-gray-400 line-clamp-3" v-html="email.content_preview"></p>
                </div>
            </div>
        </div>

        <div class="mt-6 flex justify-center items-center space-x-2">
            <span v-if="currentGroup > 0"
                @click="changeGroup(currentGroup - 1)"
                class="px-3 py-1 rounded cursor-pointer bg-[#202021] text-gray-300 hover:bg-gray-600">
                ←
            </span>

            <span v-for="page in paginationNumbers" 
                :key="page" 
                @click="changePage(page)"
                :class="(currentPage + 1) === page ? 'bg-blue-500 text-white' : 'bg-[#202021] text-gray-300 hover:bg-gray-600'"
                class="px-3 py-1 rounded cursor-pointer transition-colors">
                {{ page }}
            </span>

            <span v-if="hasNextGroup"
                @click="changeGroup(currentGroup + 1)"
                class="px-3 py-1 rounded cursor-pointer bg-[#202021] text-gray-300 hover:bg-gray-600">
                →
            </span>
        </div>
    </section>
</template>

<script>
import { ref, computed, watch } from 'vue';
import { searchEmails } from '../../services/emailService';

export default {
    name: 'EmailList',
    props: {
        selectedTab: {
            type: String,
            default: 'inbox'
        }
    },
    emits: ['emailSelected'],
    setup(props, { emit }) {
        const searchQuery = ref('');
        const fromFilter = ref('');
        const toFilter = ref('');
        const currentPage = ref(0);
        const emailsPerPage = 5;
        const totalEmails = ref(0);
        const emails = ref([]);
        const isLoading = ref(false);
        const hasSearched = ref(false);

        let debounceTimeout;
        const debounceSearch = () => {
            clearTimeout(debounceTimeout);
            debounceTimeout = setTimeout(() => {
                // Si se elimina el texto de búsqueda, limpiar los filtros
                if (!searchQuery.value.trim()) {
                    fromFilter.value = '';
                    toFilter.value = '';
                }
                resetSearch();
            }, 300);
        };

        // Función para reiniciar la búsqueda
        const resetSearch = () => {
            currentPage.value = 0;
            hasSearched.value = true;
            return onSearch(0);
        };

        // Función para manejar la búsqueda desde los inputs
        const handleSearch = () => {
            if (
                searchQuery.value.trim() === '' &&
                fromFilter.value.trim() === '' &&
                toFilter.value.trim() === ''
            )
                return;

            resetSearch();
        };

        const onSearch = async (newPage = 0) => {
            isLoading.value = true;
            try {
                const result = await searchEmails({
                    query: searchQuery.value,
                    from: fromFilter.value,
                    to: toFilter.value,
                    page: Math.floor(newPage / emailsPerPage),
                    size: emailsPerPage * 5
                });
                console.log(result);
                totalEmails.value = result.total || 0;
                emails.value = result.emails?.map((email, index) => {
                    const d = new Date(email.date);
                    const formattedDate = d.toLocaleDateString('en-US', {
                        day: '2-digit',
                        month: 'short',
                        year: 'numeric'
                    });
                    return {
                        id: index + 1,
                        sender: email.from,
                        sender_preview: email.highlight?.from ? email.highlight.from[0] : email.from,
                        receiver: email.to,
                        subject: email.subject,
                        subject_preview:
                            (email.highlight?.subject ? email.highlight.subject[0] : email.subject).length > 35
                                ? (email.highlight?.subject ? email.highlight.subject[0] : email.subject).substring(0, 35) + '...'
                                : (email.highlight?.subject ? email.highlight.subject[0] : email.subject),
                        date: formattedDate,
                        content: email.body,
                        content_preview:
                            (email.highlight?.body ? email.highlight.body[0] : email.body).length > 200
                                ? (email.highlight?.body ? email.highlight.body[0] : email.body).substring(0, 200) + '...'
                                : (email.highlight?.body ? email.highlight.body[0] : email.body)
                    };
                }) || [];
            } catch (error) {
                console.error("Error al buscar emails:", error);
            } finally {
                isLoading.value = false;
            }
        };

        watch(() => props.selectedTab, () => {
            currentPage.value = 0;
            emails.value = [];
            searchQuery.value = '';
            fromFilter.value = '';
            toFilter.value = '';
            hasSearched.value = false;
        });

        const totalPages = computed(() => Math.ceil(totalEmails.value / emailsPerPage));
        const paginatedEmails = computed(() => {
            const pageIndex = currentPage.value % 5;
            const start = pageIndex * emailsPerPage;
            return emails.value.slice(start, start + emailsPerPage);
        });

        const currentGroup = computed(() => Math.floor(currentPage.value / 5));
        
        const hasNextGroup = computed(() => {
            const lastPageInCurrentGroup = (currentGroup.value + 1) * 5;
            return lastPageInCurrentGroup < totalPages.value;
        });

        const paginationNumbers = computed(() => {
            const total = totalPages.value;
            const start = currentGroup.value * 5 + 1;
            const end = Math.min(start + 4, total);
            
            let pages = [];
            for (let i = start; i <= end; i++) {
                pages.push(i);
            }
            return pages;
        });

        const changeGroup = async (newGroup) => {
            await onSearch(newGroup * 5);
            currentPage.value = newGroup * 5;
        };

        const changePage = async (page) => {
            const newPage = page - 1;
            const newGroup = Math.floor(newPage / 5);
            
            if (newGroup !== currentGroup.value) {
                await onSearch(newGroup * 5);
            }
            
            currentPage.value = newPage;
        };

        const selectEmail = (email) => {
            emit('emailSelected', {
                email,
                searchTerms: {
                    query: searchQuery.value,
                    from: fromFilter.value,
                    to: toFilter.value
                }
            });
        };

        // Computed para mostrar/ocultar el botón X
        const hasAnyFilter = computed(() => {
            return searchQuery.value.trim() !== '' || 
                   fromFilter.value.trim() !== '' || 
                   toFilter.value.trim() !== '';
        });

        // Función para limpiar todos los filtros
        const clearFilters = () => {
            searchQuery.value = '';
            fromFilter.value = '';
            toFilter.value = '';
            hasSearched.value = false;
            emails.value = [];
            totalEmails.value = 0;
        };

        return {
            searchQuery,
            fromFilter,
            toFilter,
            currentPage,
            paginatedEmails,
            paginationNumbers,
            changePage,
            onSearch: handleSearch,
            selectEmail,
            isLoading,
            currentGroup,
            hasNextGroup,
            changeGroup,
            totalEmails,
            hasSearched,
            debounceSearch,
            hasAnyFilter,
            clearFilters
        };
    }
};
</script>

<style scoped>
:deep(mark) {
    background-color: #052f74;
    color: inherit;
    color: #fff;
    padding: 0;
}
</style>