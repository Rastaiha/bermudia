<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-gray-100 font-vazirmatn">
    <div class="bg-white shadow-lg rounded-xl p-8 w-full max-w-sm text-center">
      <h1 class="text-2xl font-bold text-cyan-700 mb-6">User Page</h1>

      <!-- User info -->
      <div class="space-y-2 mb-6">
        <p class="text-gray-700"><span class="font-semibold">Username:</span> {{ username }}</p>
        <p class="text-gray-700"><span class="font-semibold">ID:</span> {{ userId }}</p>
      </div>

      <!-- Logout Button -->
      <button
        @click="logout"
        class="w-full py-2 px-4 bg-red-500 hover:bg-red-600 text-white font-medium rounded-lg transition"
      >
        Logout
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"

const username = ref("")
const userId = ref("")
const router = useRouter()
const BASE_URL = "http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1";


onMounted(() => {
  debugger;
  let token = localStorage.getItem("authToken")
  if (!token) token = sessionStorage.getItem("authToken")

  if (!token) {
    logout();
    return
  }

  fetchUser(token);
})

async function fetchUser(token) {
  try {
    const response = await fetch(`${BASE_URL}/me`, {
      method: "GET",
      headers: {
        "Authorization": `Bearer ${token}`,
      },
    })

    const data = await response.json();

    if (!response.ok) {
      logout();
      return;
    }
    username.value = data.result.username;
    userId.value = data.result.id;
  } catch (err) {
    logout();
  }
}

function logout() {
  localStorage.removeItem("authToken");
  sessionStorage.removeItem("authToken");
  router.push("/login");
}
</script>
