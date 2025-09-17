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
    /* Adjust top and left to position it near your toolbar */
    top: 80px;
    left: 20px;
    background-color: #4caf50; /* A pleasant green */
    color: white;
    padding: 10px 15px;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    z-index: 1000; /* Ensure it's above other content */
    font-family:
        'Pelak', sans-serif; /* Assuming Pelak font is globally available */
    font-size: 1rem;
    animation: fadeIn 0.5s ease-out;
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
