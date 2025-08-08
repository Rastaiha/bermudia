<template>
  <SimpleLayout>
    <div dir="rtl" class="bg-gray-50 min-h-screen py-12">
      <div class="max-w-3xl mx-auto text-right space-y-12 px-6">
        <!-- Header Section -->
        <section class="space-y-4">
          <h1 class="text-5xl font-extrabold font-vazir text-rasta-green text-indigo-800">
            رستا
          </h1>
          <p class="text-2xl font-light font-vazir text-gray-700">
            جمع علمی‑‌ترویجی
          </p>
        </section>

        <!-- About Section -->
        <section class="space-y-4">
          <h2 class="text-3xl font-bold font-vazir text-indigo-700">
            ما کیستیم؟
          </h2>
          <p class="text-lg font-vazir text-gray-800 leading-relaxed">
            جمع علمی‑ترویجی رستا، تابستان ۱۳۹۶ توسط جمعی از دانشجویان دانشگاه‌های صنعتی شریف، صنعتی اصفهان و دانشگاه تهران با حمایت خانه ریاضیات اصفهان و دانشگاه صنعتی شریف تأسیس شد.
            رسالت رستا؛ تلاش جهت تحقق آموزش و پژوهش در راستای رشد و تربیت است. بر همین اساس، جمع علمی‑‌ترویجی رستا در سال‌های گذشته تلاش کرده با استفاده از توانایی‌های علمی اعضای خود، قدم کوچکی در راستای تقویت روحیه پرسشگری در میان دانش‌آموزان و دانشجویان کشور بردارد. در این راه رستا محتواهای آموزشی متنوعی تولید کرده و آن‌ها را در قالب رویدادهای حضوری و مجازی در اختیار دانش‌آموزان سراسر کشور قرار داده است.
          </p>
        </section>

        <!-- Events Section -->
        <section class="space-y-6">
          <h2 class="text-3xl font-bold font-vazir text-indigo-700">
            آخرین رویدادها
          </h2>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div v-for="event in events" :key="event.title" class="bg-white rounded-lg shadow-lg p-6">
              <h3 class="text-2xl font-semibold font-vazir text-gray-800 mb-2">
                {{ event.title }}
              </h3>
              <p class="text-sm font-vazir text-gray-600">{{ event.date }}</p>
              <p class="text-base font-vazir text-gray-700 mt-2">{{ event.desc }}</p>
            </div>
          </div>

          <div class="text-center mt-4">
            <button class="px-6 py-2 bg-indigo-600 text-white rounded-md font-vazir hover:bg-indigo-700 transition">
              مشاهده تمام رویدادها
            </button>
          </div>
        </section>

        <!-- Blog Posts Section -->
        <section class="space-y-6">
          <h2 class="text-3xl font-bold font-vazir text-indigo-700">
            آخرین پست‌های بلاگ
          </h2>
          <ul class="space-y-4">
            <li v-for="post in posts" :key="post.title" class="bg-white rounded-lg shadow-md p-4">
              <h3 class="text-2xl font-semibold font-vazir text-gray-800">
                {{ post.title }}
              </h3>
              <p class="text-sm font-vazir text-gray-600">{{ post.author }}</p>
              <div class="mt-2">
                <button class="text-indigo-600 font-vazir hover:underline">
                  ادامه مطلب ...
                </button>
              </div>
            </li>
          </ul>
          <div class="text-center mt-4">
            <button class="px-6 py-2 bg-indigo-600 text-white rounded-md font-vazir hover:bg-indigo-700 transition">
              مشاهده تمام پست‌ها
            </button>
          </div>
        </section>

        <!-- Footer -->
        <footer class="pt-12 text-center text-gray-500 text-sm font-vazir">
          جمع علمی‑ترویجی رستا<br />
          تمامی حقوق این سایت متعلق به جمع علمی‑ترویجی رستا است.
        </footer>
      </div>
    </div>
  </SimpleLayout>
</template>

<script setup>
import SimpleLayout from "../layouts/SimpleLayout.vue"

// Static data based on rastaiha.ir
const events = [
  {
    title: "مدرسه بهاره رستا در بافق",
    date: "اردیبهشت ۱۴۰۴",
    desc: "رویداد علمی در استان یزد شهرستان بافق"
  },
  {
    title: "رویداد علمی رستاکلا ۱۴۰۴",
    date: "اردیبهشت ۱۴۰۴",
    desc: "رویدادی علمی در دو نوبت برای معلمان و دانش‑آموزان"
  },
  {
    title: "مدرسهٔ تابستانهٔ ۱۴۰۳",
    date: "شهریور ۱۴۰۳",
    desc: "رویدادی با کارگاه‌های جدید در دو شهر تهران و اصفهان"
  }
]

const posts = [
  {
    title: "حرف‌هایی با همه در مورد اقتصاد",
    author: "سید پارسا قزوینی، دانشجوی ارشد اقتصاد در دانشگاه شریف"
  },
  {
    title: "رادیو نیم‌خط، نوروز ۱۴۰۳",
    author: "هفت‌‌سین پریسا"
  },
  {
    title: "رادیو نیم‌خط، فصل چهارم، شماره دوم",
    author: "قیچی"
  }
]
</script>
