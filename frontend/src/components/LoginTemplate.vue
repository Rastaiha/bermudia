<template>
  <div class="relative min-h-screen flex items-center justify-center bg-sky-400 overflow-hidden font-vazirmatn">
    <!-- Waves Back -->
    <div class="waves back">
      <svg viewBox="0 24 150 28" preserveAspectRatio="none" class="w-full h-full">
        <defs>
          <path
            id="gentle-wave"
            d="M-160 44c30 0 58-18 88-18s 58 18 88 18
               58-18 88-18 58 18 88 18 v44h-352z"
          />
        </defs>
        <use href="#gentle-wave" x="48" y="5" />
        <use href="#gentle-wave" x="48" y="7" />
      </svg>
    </div>

    <!-- Boat -->
    <div class="boat">
      <div class="sail2"></div>
      <div class="sail1"></div>
      <div class="hull"></div>
    </div>

    <!-- Waves Front -->
    <div class="waves front">
      <svg viewBox="0 24 150 28" preserveAspectRatio="none" class="w-full h-full">
        <use href="#gentle-wave" x="48" y="0" />
        <use href="#gentle-wave" x="48" y="3" />
      </svg>
    </div>

    <!-- Login Card -->
    <div class="animate-fade-in-scale w-full max-w-md mx-4 relative z-10">
      <div class="card-hover bg-white rounded-xl shadow-lg overflow-hidden border border-blue-100">
        <div class="p-6 sm:p-8">
          <div class="flex justify-center mb-8">
            <div class="bg-cyan-100 p-4 rounded-full">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-cyan-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </div>
          </div>

          <h1 class="text-2xl font-bold text-center text-cyan-800 mb-1">ورود به سیستم</h1>

          <form @submit.prevent="handleLogin" class="space-y-6">
            <!-- username -->
            <div>
              <label for="username" class="block text-sm font-medium text-cyan-700 mb-1">نام کاربری</label>
              <div class="relative">
                <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-cyan-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                </div>
                <input
                  id="username"
                  v-model="username"
                  autocomplete="on"
                  type="text"
                  class="input-focus text-left bg-blue-50 border border-cyan-200 text-cyan-800 text-sm sm:text-base rounded-lg block w-full pl-10 p-2 sm:p-3 sm:pl-10"
                  placeholder="Username"
                  required
                />
              </div>
            </div>

            <!-- Password -->
            <div>
              <label for="password" class="block text-sm font-medium text-cyan-700 mb-1">رمز عبور</label>
              <div class="relative">
                <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-cyan-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                  </svg>
                </div>
                <input
                  id="password"
                  v-model="password"
                  type="password"
                  autocomplete="on"
                  class="input-focus text-left bg-blue-50 border border-cyan-200 text-cyan-800 text-sm sm:text-base rounded-lg block w-full pl-10 p-2 sm:p-3 sm:pl-10"
                  placeholder="********"
                  required
                />
              </div>
            </div>

            <!-- Remember -->
            <div class="flex items-center justify-between">
              <div class="flex items-start">
                <div class="flex items-center h-5">
                  <input id="remember" type="checkbox"
                   class="w-4 h-4 bg-blue-50 rounded border border-cyan-300 focus:ring-3 focus:ring-cyan-200" v-model="remember"
                  >
                </div>
                <div class="mr-3 text-sm">
                  <label for="remember" class="font-medium text-cyan-700">مرا به خاطر بسپار</label>
                </div>
              </div>
              <div></div>
            </div>

            <!-- Button -->
            <button type="submit"
              class="btn-hover w-full text-white bg-cyan-600 hover:bg-cyan-700 font-medium rounded-lg text-sm px-5 py-2.5 text-center">
              لاگین
            </button>

            <p v-if="error" class="bg-red-500 text-white text-sm text-center mt-3 rounded-lg px-5 py-2.5">{{ error }}</p>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from "vue";
import { useTimeout } from "./service/ChainedTimeout.js"; // import your timeout service

const username = ref("");
const password = ref("");
const error = ref("");
const token = ref("");
const boatLeft = ref(10);
let boatEl;
const remember = ref(false);

const { startTimeout, clear } = useTimeout(); // initialize timeout

const BASE_URL = "http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1";

async function handleLogin() {
  try {
    const response = await fetch(`${BASE_URL}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username: username.value, password: password.value }),
    });

    const data = await response.json();

    if (!response.ok) {
      error.value = data.error;
      startTimeout(() => {
        error.value = "";
      }, 5000);

      return;
    }

    token.value = data.token;
    error.value = "";
    if (remember.value) {
      localStorage.setItem("authToken", token.value);
    }
    window.location.href = "/user_page";
  } catch (err) {
    error.value = err.message;
    startTimeout(() => {
      error.value = "";
    }, 5000);
  }
}

function moveBoat(e) {
  if (!boatEl) return;
  if (e.key === "ArrowLeft" && boatLeft.value > 5) boatLeft.value -= 5;
  if (e.key === "ArrowRight" && boatLeft.value < 60) boatLeft.value += 5;
  boatEl.style.left = `${boatLeft.value}vw`;
}

onMounted(() => {
  boatEl = document.querySelector(".boat");
  window.addEventListener("keydown", moveBoat);
});
onBeforeUnmount(() => {
  window.removeEventListener("keydown", moveBoat);
  clear(); // clear any pending timeout on unmount
});
</script>

<style scoped>
/* Animations */
@keyframes fadeInScale {
  0% { opacity: 0; transform: scale(0.95); }
  100% { opacity: 1; transform: scale(1); }
}
.animate-fade-in-scale { animation: fadeInScale 0.6s ease-out forwards; }

.card-hover { transition: all 0.3s ease, transform 0.3s ease; }
.card-hover:hover { transform: scale(1.02); box-shadow: 0 10px 25px -5px rgba(0, 105, 146, 0.2); }

.input-focus { transition: all 0.3s ease; }
.input-focus:focus {
  border-color: #38bdf8;
  box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.2);
}

.btn-hover { transition: all 0.3s ease, transform 0.2s ease; }
.btn-hover:hover { background-color: #0ea5e9; transform: scale(1.03); }

/* Waves */
.waves { position: fixed; bottom: 50px; left: 0; height: 50vh; width: 100%; }
.waves::after {
  content: ""; display: block; position: absolute; left: 0; width: 100%;
  bottom: -100px; height: 100px; background: dodgerblue;
}
.waves svg { position: absolute; bottom: 0; left: 0; }
.waves use { animation: wavewave 5s cubic-bezier(.55,.5,.45,.5) infinite; fill: dodgerblue; opacity: 0.7; }
.waves.back { z-index: 1; }
.waves.back use:nth-child(1) { animation-delay: -4s; animation-duration: 13s; }
.waves.back use:nth-child(2) { animation-delay: -5s; animation-duration: 20s; }
.waves.front { z-index: 3; }
.waves.front use:nth-child(1) { animation-delay: -2s; animation-duration: 7s; }
.waves.front use:nth-child(2) { animation-delay: -3s; animation-duration: 10s; }
@keyframes wavewave { 0% { transform: translate(-90px); } 100% { transform: translate(85px); } }

/* Boat */
.boat { width: 25vw; position: fixed; bottom: 50px; left: 10vw; z-index: 2; display: inline-block;
  transform-origin: center bottom; animation: bob 4s ease-in-out alternate infinite; }
.boat .hull { width: 100%; height: calc(25vw / 20 * 7); background: firebrick;
  border-radius: 0 0 12.5vw 12.5vw; border-top: 0.375vw solid darkred; }
.boat .sail1 { width: 0; border-bottom: calc(25vw / 20 * 9) solid white; border-left: calc(25vw / 20 * 8) solid transparent;
  margin: 0 calc(25vw / 40) calc(25vw / 100); display: inline-block; }
.boat .sail2 { width: 0; border-bottom: calc(25vw / 20 * 11) solid white; border-right: calc(25vw / 2) solid transparent; display: inline-block; }
@keyframes bob { 0% { transform: rotate(-10deg) translateY(0vh); } 100% { transform: rotate(8deg) translateY(-20vh); } }
</style>
