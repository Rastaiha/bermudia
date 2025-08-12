import { onBeforeUnmount } from 'vue'

export function useTimeout() {
  let timeoutId = null

  const startTimeout = (callback, delay) => {
    clearTimeout(timeoutId) // clear previous timeout if any
    timeoutId = setTimeout(callback, delay)
  }

  const clear = () => {
    clearTimeout(timeoutId)
    timeoutId = null
  }

  onBeforeUnmount(() => {
    clear()
  })

  return { startTimeout, clear }
}