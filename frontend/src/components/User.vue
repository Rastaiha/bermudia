<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-gray-100 font-vazirmatn">
    <div class="bg-white shadow-lg rounded-xl p-8 w-full max-w-sm text-center">
      <h1 class="text-2xl font-bold text-cyan-700 mb-6">صفحه کاربر</h1>

      <div class="space-y-2 mb-6">
        <p class="text-gray-700"><span class="font-semibold">نام کاربری:</span> {{ username }}</p>
        <p class="text-gray-700"><span class="font-semibold">شناسه:</span> {{ userId }}</p>
      </div>
      
      <p v-if="errorMsg" class="bg-red-100 text-red-700 text-sm text-center my-3 rounded-lg px-5 py-2.5">{{ errorMsg }}</p>

      <div class="border-t pt-6 space-y-3">
        <button
          @click="enterTerritory"
          :disabled="isLoading"
          class="w-full py-2.5 px-4 bg-cyan-600 hover:bg-cyan-700 text-white font-medium rounded-lg transition disabled:bg-gray-400"
        >
          <span v-if="!isLoading">ورود به جزیره</span>
          <span v-else>در حال بارگذاری...</span>
        </button>

        <button
          @click="handleLogout"
          class="w-full py-2 px-4 bg-red-500 hover:bg-red-600 text-white font-medium rounded-lg transition"
        >
          خروج
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
// Import from the API service
import { getMe, getPlayer, logout } from "@/services/api";

const username = ref("...");
const userId = ref("...");
const errorMsg = ref("");
const isLoading = ref(false);
const router = useRouter();

// Fetch user data on component mount
onMounted(async () => {
  try {
    const meData = await getMe();
    username.value = meData.username;
    userId.value = meData.id;
  } catch (err) {
    // If getMe fails, token is likely invalid or expired
    errorMsg.value = "نشست شما منقضی شده است. لطفاً دوباره وارد شوید.";
    setTimeout(handleLogout, 2000);
  }
});

// Function to get player data and navigate to their territory
async function enterTerritory() {
  isLoading.value = true;
  errorMsg.value = "";
  try {
    const playerData = await getPlayer();
    if (playerData && playerData.atTerritory) {
      router.push(`/territory/${playerData.atTerritory}`);
    } else {
      throw new Error("اطلاعات تریتوری از سرور دریافت نشد.");
    }
  } catch (err) {
    console.error("Failed to get player data:", err);
    errorMsg.value = err.message;
  } finally {
    isLoading.value = false;
  }
}

// Function to log out the user
function handleLogout() {
  logout();
  router.push("/login");
}
</script>