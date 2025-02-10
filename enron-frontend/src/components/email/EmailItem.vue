<template>
    <div 
        class="flex items-center p-4 bg-[#202021] rounded-lg shadow-md cursor-pointer hover:bg-[#202021] transition-colors h-32"
        @click="$emit('select', email)"
    >
        <div class="flex-shrink-0 w-10 h-10 rounded-full bg-gray-600 flex items-center justify-center text-white font-bold mr-4">
            {{ email.sender.charAt(0).toUpperCase() }}
        </div>
        <div class="flex-1">
            <div class="flex justify-between items-center">
                <span class="font-bold truncate" v-html="email.sender_preview"></span>
                <span class="text-sm text-gray-400">{{ email.date }}</span>
            </div>
            <h3 class="text-m font-semibold truncate" v-html="email.subject_preview"></h3>
            <p class="text-sm text-gray-400 line-clamp-3" v-html="email.content_preview"></p>
        </div>
    </div>
</template>

<script>
export default {
    name: 'EmailItem',
    props: {
        email: {
            type: Object,
            required: true
        }
    },
    emits: ['select']
}
</script> 