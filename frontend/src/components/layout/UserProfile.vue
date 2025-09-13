<template>
    <div dir="rtl" class="fixed top-4 right-4 z-50">
        <div
            class="profile-container flex items-center rounded-full bg-black/20 backdrop-blur-sm transition-all duration-300 ease-in-out"
        >
            <div
                class="profile-image w-14 h-14 rounded-full bg-gray-200 overflow-hidden flex items-center justify-center"
            >
                <img
                    :src="profileImageUrl"
                    alt="User Profile Image"
                    class="w-full h-full object-cover"
                />
            </div>
            <span class="text-white font-medium text-lg pr-4 pl-5 py-2">{{
                username
            }}</span>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
    username: {
        type: String,
        required: true,
    },
});

const getProfileImageNumber = username => {
    if (!username) return 1;
    let sum = 0;
    for (let i = 0; i < username.length; i++) {
        sum += username.charCodeAt(i);
    }
    return (sum % 15) + 1;
};

const profileImageUrl = computed(() => {
    const imageNumber = getProfileImageNumber(props.username);
    return `/images/profiles/${imageNumber}.png`;
});
</script>

<style scoped>
.profile-container {
    transition: all 0.3s ease;
}
.profile-container:hover .profile-image {
    transform: scale(1.1);
}
.profile-image {
    transition: transform 0.3s ease;
}
</style>
