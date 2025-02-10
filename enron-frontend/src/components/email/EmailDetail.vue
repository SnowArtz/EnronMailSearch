<template>
    <aside class="w-[65%] bg-white p-6 shadow-xl rounded-xl mx-2 my-2">
        <div v-if="email"
            class="bg-white p-8 rounded-lg max-h-[90vh] overflow-y-auto scrollbar-thin scrollbar-thumb-gray-400 scrollbar-track-gray-100">
            <div class="flex justify-between items-start mb-6">
                <div class="flex items-center space-x-4">
                    <div class="w-12 h-12 rounded-full bg-gray-600 flex items-center justify-center text-white font-bold">
                        {{ email.sender.charAt(0).toUpperCase() }}
                    </div>
                    <div class="flex-1">
                        <h2 class="text-2xl font-bold truncate" v-html="highlightedSubject"></h2>
                        <div class="flex flex-col text-sm text-gray-500">
                            <span>De: <span v-html="highlightedSender"></span></span>
                            <span>Para: <span v-html="highlightedReceiver"></span></span>
                            <span>{{ formattedDate }}</span>
                        </div>
                    </div>
                </div>
                <BaseButton 
                    @click="$emit('closeEmail')" 
                    variant="icon"
                    title="Cerrar correo"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </BaseButton>
            </div>
            <p class="text-gray-700 whitespace-pre-wrap" v-html="highlightedContent"></p>
        </div>
        <div v-else class="flex flex-col items-center justify-center h-full">
            <svg class="w-24 h-24 text-gray-500 mb-4" fill="none" stroke="currentColor" stroke-width="1" viewBox="0 0 24 24">
                <path d="M22 6c0-1.1-.9-2-2-2H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6zm-2 0l-8 5-8-5h16zm0 12H4V8l8 5 8-5v10z"/>
            </svg>
            <span class="text-gray-500 text-lg">Selecciona un elemento para leerlo</span>
        </div>
    </aside>
</template>

<script>
import { computed } from 'vue';
import { highlightText } from '../../utils/textHighlighter';
import { formatDate } from '../../utils/dateFormatter';
import BaseButton from '../ui/BaseButton.vue';

export default {
    name: 'EmailDetail',
    components: {
        BaseButton
    },
    props: {
        email: {
            type: Object,
            default: null
        },
        searchTerms: {
            type: Object,
            default: () => ({
                query: '',
                from: '',
                to: ''
            })
        }
    },
    emits: ['closeEmail'],
    setup(props) {
        const highlightedSubject = computed(() => 
            highlightText(props.email?.subject, props.searchTerms)
        );

        const highlightedSender = computed(() => 
            highlightText(props.email?.sender, props.searchTerms)
        );

        const highlightedReceiver = computed(() => 
            highlightText(props.email?.receiver, props.searchTerms)
        );

        const highlightedContent = computed(() => 
            highlightText(props.email?.content, props.searchTerms)
        );

        const formattedDate = computed(() => 
            props.email ? formatDate(props.email.date) : ''
        );

        return {
            highlightedSubject,
            highlightedSender,
            highlightedReceiver,
            highlightedContent,
            formattedDate
        };
    }
};
</script>

<style scoped>
:deep(mark) {
    background-color: #084dbd;
    color: inherit;
    color: #fff;
    padding: 2px;
    border-radius: 3px;
    opacity: 0.8;
}
</style>