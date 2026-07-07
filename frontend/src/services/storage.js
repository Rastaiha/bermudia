export const getJSON = (key, fallback = []) => {
    try {
        const data = localStorage.getItem(key);
        return data ? JSON.parse(data) : fallback;
    } catch (e) {
        console.error(`Error reading ${key} from localStorage:`, e);
        return fallback;
    }
};

export const setJSON = (key, value) => {
    try {
        localStorage.setItem(key, JSON.stringify(value));
    } catch (e) {
        console.error(`Error saving ${key} to localStorage:`, e);
    }
};
