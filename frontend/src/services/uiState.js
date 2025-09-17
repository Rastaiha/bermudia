import { reactive } from 'vue';

// A global reactive state to track UI properties
export const uiState = reactive({
    isInboxOpen: false,
});
