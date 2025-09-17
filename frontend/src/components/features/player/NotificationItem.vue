<template>
    <div
        class="bg-gray-700 rounded-xl shadow-md overflow-hidden transition-all duration-300 hover:shadow-lg hover:bg-gray-600"
        :class="notification.style.border"
    >
        <div class="p-4">
            <div class="flex items-start justify-between">
                <div class="flex items-start gap-4">
                    <div
                        class="p-2 rounded-full flex items-center justify-center"
                        :class="notification.style.bg"
                    >
                        <component
                            :is="notification.icon"
                            :class="notification.style.icon"
                            class="w-7 h-7"
                        />
                    </div>
                    <div>
                        <h3 class="font-bold text-gray-100">
                            {{ notification.title }}
                        </h3>
                        <p
                            class="text-gray-300 text-sm mt-1"
                            v-html="notification.summary"
                        ></p>
                    </div>
                </div>
                <span class="text-xs text-gray-400 flex-shrink-0">{{
                    formattedDate
                }}</span>
            </div>
            <div v-if="notification.details" class="mt-3">
                <button
                    class="text-sm font-medium flex items-center transition-colors"
                    :class="notification.style.text"
                    @click="toggleExpand"
                >
                    جزئیات بیشتر
                    <ChevronDownIcon
                        class="w-4 h-4 mr-1 transition-transform duration-300"
                        :class="{ 'rotate-180': isExpanded }"
                    />
                </button>
                <div
                    class="expand-content mt-2"
                    :class="{ expanded: isExpanded }"
                >
                    <div
                        class="p-4 rounded-lg text-sm"
                        :class="notification.style.detailsBg"
                    >
                        <div
                            class="text-gray-200 leading-relaxed"
                            v-html="notification.details"
                        ></div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import {
    CheckCircleIcon,
    XCircleIcon,
    InformationCircleIcon,
    CurrencyDollarIcon,
    ChevronDownIcon,
} from '@heroicons/vue/24/outline';
import { formatDistanceToNow } from 'date-fns-jalali';
import { COST_ITEMS_INFO } from '@/services/cost.js';
import { glossary } from '@/services/glossary.js';

const props = defineProps({
    message: {
        type: Object,
        required: true,
    },
});

const isExpanded = ref(false);

const toggleExpand = () => {
    isExpanded.value = !isExpanded.value;
};

const formattedDate = computed(() => {
    if (!props.message.createdAt) return '';
    return formatDistanceToNow(new Date(parseInt(props.message.createdAt)), {
        addSuffix: true,
    });
});

const formatItemsToList = items => {
    if (!items || items.length === 0) return '';
    const listItems = items
        .map(item => {
            const details = COST_ITEMS_INFO[item.type] || {
                name: item.type,
                icon: '',
            };
            return `<div>- ${details.name} (<img src="${details.icon}" class="w-4 h-4 inline-block mx-0.5" />x${item.amount})</div>`;
        })
        .join('');
    return `<div class="mt-2">${listItems}</div>`;
};

const formatItemsToGrid = items => {
    if (!items || items.length === 0) return '<div>-</div>';
    return items
        .map(item => {
            const details = COST_ITEMS_INFO[item.type] || {
                name: item.type,
                icon: '',
            };
            return `<div class="flex items-center gap-2" title="${details.name}">
                      <img src="${details.icon}" class="w-6 h-6" />
                      <span class="font-semibold text-base text-gray-200">${item.amount}</span>
                    </div>`;
        })
        .join('');
};

const notification = computed(() => {
    const content = props.message.content;
    if (content.newCorrection) {
        const correction = content.newCorrection;
        const status = correction.newState.status;

        let title, icon, style, statusText;

        switch (status) {
            case 'correct':
                title = 'پاسخ شما صحیح بود';
                icon = CheckCircleIcon;
                statusText = 'صحیح';
                style = {
                    border: 'border-r-4 border-green-500',
                    bg: 'bg-green-800',
                    icon: 'text-green-300',
                    text: 'text-green-400 hover:text-green-300',
                    detailsBg: 'bg-green-900/50',
                };
                break;
            case 'half-correct':
                title = 'پاسخ شما نیمه‌درست بود';
                icon = InformationCircleIcon;
                statusText = 'نیمه‌درست';
                style = {
                    border: 'border-r-4 border-amber-500',
                    bg: 'bg-amber-800',
                    icon: 'text-amber-300',
                    text: 'text-amber-400 hover:text-amber-300',
                    detailsBg: 'bg-amber-900/50',
                };
                break;
            default: // wrong
                title = 'پاسخ شما صحیح نبود';
                icon = XCircleIcon;
                statusText = 'غلط';
                style = {
                    border: 'border-r-4 border-red-500',
                    bg: 'bg-red-800',
                    icon: 'text-red-300',
                    text: 'text-red-400 hover:text-red-300',
                    detailsBg: 'bg-red-900/50',
                };
                break;
        }

        let details = `پاسخ شما برای سوال در ${glossary.island} <strong>${correction.islandName || 'نامشخص'}</strong> در ${glossary.territory} <strong>${correction.territoryName || 'نامشخص'}</strong> تصحیح شد. <br/> وضعیت: <strong>${statusText}</strong>.`;

        if (correction.newState.feedback) {
            details += `<div class="my-2.5 border-t border-gray-600"></div><span>بازخورد:</span><div class="mt-1 text-justify text-gray-300">${correction.newState.feedback}</div>`;
        }

        if (correction.reward?.items?.length > 0) {
            const rewardList = formatItemsToList(correction.reward.items);
            details += `<div class="my-2.5"></div><span>جایزه شما:</span>${rewardList}`;
        }

        return {
            title,
            summary: `در ${glossary.island} «${correction.islandName}»`,
            details: details,
            icon,
            style,
        };
    }
    if (content.announcement) {
        return {
            title: 'اعلامیه',
            summary: content.announcement.text,
            details: null,
            icon: InformationCircleIcon,
            style: {
                border: 'border-r-4 border-teal-500',
                bg: 'bg-teal-800',
                icon: 'text-teal-300',
                text: 'text-teal-400 hover:text-teal-300',
                detailsBg: 'bg-teal-900/50',
            },
        };
    }
    if (content.ownOfferAccepted) {
        const offer = content.ownOfferAccepted.offer;
        const offeredGrid = formatItemsToGrid(offer.offered.items);
        const requestedGrid = formatItemsToGrid(offer.requested.items);

        const details = `
            <p>پیشنهاد معامله شما توسط بازیکن دیگر پذیرفته شد:</p>
            <div class="my-3 flex items-center justify-center text-center bg-gray-900/50 p-3 rounded-lg">
                <div class="flex-1 flex flex-col gap-2.5 items-center">
                    <span class="text-xs text-gray-400 mb-1">شما دادید</span>
                    ${offeredGrid}
                </div>
                <div class="flex-shrink-0">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-7 h-7 text-gray-400 mx-3">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 21 3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5" />
                    </svg>
                </div>
                <div class="flex-1 flex flex-col gap-2.5 items-center">
                    <span class="text-xs text-gray-400 mb-1">شما گرفتید</span>
                    ${requestedGrid}
                </div>
            </div>
            <p>آیتم‌های جدید به دارایی شما اضافه شد.</p>
        `;

        return {
            title: 'معامله انجام شد',
            summary: 'پیشنهاد شما در بازار پذیرفته شد.',
            details: details,
            icon: CurrencyDollarIcon,
            style: {
                border: 'border-r-4 border-blue-500',
                bg: 'bg-blue-800',
                icon: 'text-blue-300',
                text: 'text-blue-400 hover:text-blue-300',
                detailsBg: 'bg-blue-900/50',
            },
        };
    }
    return {
        title: 'پیام نامشخص',
        summary: 'یک پیام جدید دریافت کرده‌اید.',
        icon: 'div',
        style: {
            border: 'border-r-4 border-gray-500',
            bg: 'bg-gray-800',
            icon: 'text-gray-300',
            text: 'text-gray-400 hover:text-gray-300',
            detailsBg: 'bg-gray-900/50',
        },
    };
});
</script>

<style scoped>
.expand-content {
    max-height: 0;
    overflow: hidden;
    transition: max-height 0.3s ease-out;
}
.expand-content.expanded {
    max-height: 500px;
    transition: max-height 0.5s ease-in;
}
</style>
