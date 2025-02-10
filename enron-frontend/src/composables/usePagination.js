import { ref, computed } from 'vue';

export function usePagination(totalItems, itemsPerPage = 5) {
    const currentPage = ref(0);
    
    const totalPages = computed(() => 
        Math.ceil(totalItems.value / itemsPerPage)
    );

    const currentGroup = computed(() => 
        Math.floor(currentPage.value / 5)
    );

    const hasNextGroup = computed(() => {
        const lastPageInCurrentGroup = (currentGroup.value + 1) * 5;
        return lastPageInCurrentGroup < totalPages.value;
    });

    const paginationNumbers = computed(() => {
        const total = totalPages.value;
        const start = currentGroup.value * 5 + 1;
        const end = Math.min(start + 4, total);
        
        return Array.from(
            { length: end - start + 1 }, 
            (_, i) => start + i
        );
    });

    return {
        currentPage,
        totalPages,
        currentGroup,
        hasNextGroup,
        paginationNumbers
    };
} 