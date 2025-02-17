<template>
  <aside class="w-16 bg-[#202021] flex flex-col items-center py-6 space-y-6">
    <button @click="updateTab('inbox')" :class="selectedTab === 'inbox' ? 'bg-gray-600 text-white' : 'text-gray-400'"
      class="p-3 rounded-full hover:bg-gray-600 transition-colors" title="Correos">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
        class="size-6">
        <path stroke-linecap="round" stroke-linejoin="round"
          d="M21.75 6.75v10.5a2.25 2.25 0 0 1-2.25 2.25h-15a2.25 2.25 0 0 1-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0 0 19.5 4.5h-15a2.25 2.25 0 0 0-2.25 2.25m19.5 0v.243a2.25 2.25 0 0 1-1.07 1.916l-7.5 4.615a2.25 2.25 0 0 1-2.36 0L3.32 8.91a2.25 2.25 0 0 1-1.07-1.916V6.75" />
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