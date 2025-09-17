<template>
    <div v-if="hasUnread" class="notification-indicator">
        <p>شما پیام جدیدی دارید!</p>
    </div>
</template>

<script setup>
import { notificationService } from '@/services/notificationService';

const hasUnread = notificationService.hasUnreadMessages;
</script>

<style scoped>
.notification-indicator {
    position: fixed;
    top: 80px;
    left: 20px;
    background-color: #4b5563; /* gray-600 */
    color: white;
    padding: 12px 16px;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    z-index: 1000;
    font-family: 'Pelak', sans-serif;
    font-size: 1rem;
    animation: fadeIn 0.5s ease-out;
    text-align: center;
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(-20px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

p {
    margin: 0;
}
</style>
