export const logger = {
    log: (...args) => {
        if (import.meta.env.DEV) console.log(...args);
    },
    error: (...args) => console.error(...args),
    warn: (...args) => console.warn(...args),
};
