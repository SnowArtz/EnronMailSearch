<template>
    <aside class="w-16 bg-[#202021] flex flex-col items-center py-6 space-y-6">
      <button
        @click="updateTab('inbox')"
        :class="selectedTab === 'inbox' ? 'bg-gray-600 text-white' : 'text-gray-400'"
        class="p-3 rounded-full hover:bg-gray-600 transition-colors"
        title="Correos"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" 
             viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
        </svg>
      </button>
    </aside>
  </template>
  
  <script>
  export default {
    name: 'Sidebar',
    props: {
      selectedTab: {
        type: String,
        default: 'inbox'
      }
    },
    emits: ['updateTab'],
    methods: {
      updateTab(tab) {
        this.$emit('updateTab', tab);
      }
    }
  };
  </script>