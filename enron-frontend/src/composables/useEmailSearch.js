import { ref, computed } from 'vue';
import { searchEmails } from '../services/emailService';

export function useEmailSearch() {
    const searchQuery = ref('');
    const fromFilter = ref('');
    const toFilter = ref('');
    const emails = ref([]);
    const isLoading = ref(false);
    const totalEmails = ref(0);
    const hasSearched = ref(false);

    const hasAnyFilter = computed(() => {
        return searchQuery.value.trim() !== '' || 
               fromFilter.value.trim() !== '' || 
               toFilter.value.trim() !== '';
    });

    const clearFilters = () => {
        searchQuery.value = '';
        fromFilter.value = '';
        toFilter.value = '';
        hasSearched.value = false;
        emails.value = [];
        totalEmails.value = 0;
    };

    const performSearch = async (page, size) => {
        isLoading.value = true;
        try {
            const result = await searchEmails({
                query: searchQuery.value,
                from: fromFilter.value,
                to: toFilter.value,
                page,
                size
            });
            
            totalEmails.value = result.total || 0;
            emails.value = formatEmails(result.emails);
        } catch (error) {
            console.error("Error al buscar emails:", error);
        } finally {
            isLoading.value = false;
        }
    };

    return {
        searchQuery,
        fromFilter,
        toFilter,
        emails,
        isLoading,
        totalEmails,
        hasSearched,
        hasAnyFilter,
        clearFilters,
        performSearch
    };
} 