<template>
  <div class="h-screen p-4 bg-white flex items-center justify-center">
    <div class="flex h-full w-[90%] bg-[#202021] rounded-xl overflow-hidden shadow-lg">
      <!-- Componente Sidebar -->
      <Sidebar :selectedTab="selectedTab" @updateTab="selectedTab = $event" />

      <div class="flex bg-[#191819] flex-1">
        <!-- Componente EmailList -->
        <EmailList 
          ref="emailListRef"
          :selectedTab="selectedTab" 
          @emailSelected="handleEmailSelected" 
        />
        <!-- Componente EmailDetail -->
        <EmailDetail 
          :email="selectedEmail?.email" 
          :searchTerms="selectedEmail?.searchTerms"
          @closeEmail="handleCloseEmail"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import Sidebar from './components/layout/Sidebar.vue';
import EmailList from './components/email/EmailList.vue';
import EmailDetail from './components/email/EmailDetail.vue';

export default {
  name: 'App',
  components: { Sidebar, EmailList, EmailDetail },
  setup() {
    const selectedTab = ref('inbox');
    const selectedEmail = ref(null);

    const emailListRef = ref(null);

    const handleEmailSelected = (email) => {
      selectedEmail.value = email;
    };

    const handleCloseEmail = () => {
      selectedEmail.value = null;
      if (emailListRef.value && emailListRef.value.clearSelectedEmail) {
        emailListRef.value.clearSelectedEmail();
      }
    };

    return {
      selectedTab,
      selectedEmail,
      emailListRef,
      handleEmailSelected,
      handleCloseEmail,
    };
  }
};
</script>
