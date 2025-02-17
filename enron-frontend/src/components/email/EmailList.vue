<template>
    <section class="w-[45%] bg-[#191819] text-gray-300 p-6 rounded-xl flex flex-col">
        <!-- Filtros -->
        <div class="flex flex-col gap-4 mb-4">
            <!-- Barra de búsqueda -->
            <div class="relative">
                <span class="absolute inset-y-0 left-0 flex items-center pl-3">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500" fill="none"
                        viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                </span>
                <input type="text" v-model="searchQuery" @input="debounceSearch" placeholder="Buscar emails..."
                    class="w-full pl-10 pr-10 py-3 rounded-full bg-[#202021] placeholder-gray-500 text-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                <!-- X -->
                <button v-if="searchQuery.trim()" @click="clearFilters"
                    class="absolute inset-y-0 right-0 flex items-center pr-3">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500 hover:text-gray-300"
                        fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>
            <!-- Remitente y Destinatario -->
            <div v-if="searchQuery.trim()" class="flex gap-4">
                <input type="text" v-model="fromFilter" @input="debounceSearch" placeholder="Remitente"
                    class="w-full pl-4 pr-4 py-2 rounded-full bg-[#202021] placeholder-gray-500 text-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                <input type="text" v-model="toFilter" @input="debounceSearch" placeholder="Destinatario"
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

        <div :class="['space-y-4 flex-1 min-h-[600px]', { 'invisible': isLoading }]">
            <EmailItem v-for="email in paginatedEmails" :key="email.id" :email="email"
                :selected="selectedEmailId === email.id" @select="selectEmail" />
        </div>

        <div class="mt-3 flex justify-center items-center space-x-2">
            <span v-if="currentGroup > 0" @click="changeGroup(currentGroup - 1)"
                class="px-1 py-1 rounded cursor-pointer bg-[#202021] text-gray-300 hover:bg-gray-600 flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="0.5"
                    stroke="currentColor" class="size-6">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5 8.25 12l7.5-7.5" />
                </svg>
            </span>

            <span v-for="page in paginationNumbers" :key="page" @click="changePage(page)"
                :class="(currentPage + 1) === page ? 'bg-blue-500 text-white' : 'bg-[#202021] text-gray-300 hover:bg-gray-600'"
                class="px-3 py-1 rounded cursor-pointer transition-colors">
                {{ page }}
            </span>

            <span v-if="hasNextGroup" @click="changeGroup(currentGroup + 1)"
                class="px-1 py-1 rounded cursor-pointer bg-[#202021] text-gray-300 hover:bg-gray-600 flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="0.5"
                    stroke="currentColor" class="size-6">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m8.25 4.5 7.5 7.5-7.5 7.5" />
                </svg>
            </span>
        </div>
    </section>
</template>

<script>
import { ref, computed, watch } from 'vue';
import { searchEmails } from '../../services/emailService';
import { formatDate } from '../../utils/dateFormatter';
import EmailItem from './EmailItem.vue';

export default {
    name: 'EmailList',
    components: {
        EmailItem
    },
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
        const pagesPerGroup = 5;
        const totalEmails = ref(0);
        const emails = ref([]);
        const isLoading = ref(false);
        const hasSearched = ref(false);
        const selectedEmailId = ref(null);

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

        const onSearch = async (newPage = 0) => {
            isLoading.value = true;
            try {
                const result = await searchEmails({
                    query: searchQuery.value,
                    from: fromFilter.value,
                    to: toFilter.value,
                    group: Math.floor(newPage / pagesPerGroup),
                    size: emailsPerPage * pagesPerGroup
                });
                totalEmails.value = result.total || 0;
                emails.value = result.emails?.map((email) => {
                    const formattedDate = formatDate(email.date);
                    return {
                        id: email.id,
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
            const pageIndex = currentPage.value % pagesPerGroup;
            const start = pageIndex * emailsPerPage;
            return emails.value.slice(start, start + emailsPerPage);
        });

        const currentGroup = computed(() => Math.floor(currentPage.value / pagesPerGroup));

        const hasNextGroup = computed(() => {
            const lastPageInCurrentGroup = (currentGroup.value + 1) * pagesPerGroup;
            return lastPageInCurrentGroup < totalPages.value;
        });

        const paginationNumbers = computed(() => {
            const total = totalPages.value;
            const start = currentGroup.value * pagesPerGroup + 1;
            const end = Math.min(start + pagesPerGroup - 1, total);

            let pages = [];
            for (let i = start; i <= end; i++) {
                pages.push(i);
            }
            return pages;
        });

        const changeGroup = async (newGroup) => {
            await onSearch(newGroup * pagesPerGroup);
            currentPage.value = newGroup * pagesPerGroup;
        };

        const changePage = async (page) => {
            const newPage = page - 1;
            const newGroup = Math.floor(newPage / pagesPerGroup);

            if (newGroup !== currentGroup.value) {
                await onSearch(newGroup * pagesPerGroup);
            }

            currentPage.value = newPage;
        };

        const selectEmail = (email) => {
            selectedEmailId.value = email.id;
            emit('emailSelected', {
                email,
                searchTerms: {
                    query: searchQuery.value,
                    from: fromFilter.value,
                    to: toFilter.value
                }
            });
        };

        const clearFilters = () => {
            searchQuery.value = '';
            fromFilter.value = '';
            toFilter.value = '';
            hasSearched.value = false;
            emails.value = [];
            totalEmails.value = 0;
            currentPage.value = 0;
        };

        const clearSelectedEmail = () => {
            selectedEmailId.value = null;
        };

        return {
            searchQuery,
            fromFilter,
            toFilter,
            currentPage,
            paginatedEmails,
            paginationNumbers,
            changePage,
            selectEmail,
            isLoading,
            currentGroup,
            hasNextGroup,
            changeGroup,
            totalEmails,
            hasSearched,
            debounceSearch,
            clearFilters,
            clearSelectedEmail,
            selectedEmailId
        };
    }
};
</script>

<style scoped>
:deep(mark) {
    background-color: #052f74;
    color: #fff;
    padding: 2px;
    border-radius: 5px;
}
</style>