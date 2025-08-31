// frontend\src\components\service\ChainedTimeout.js

import { onBeforeUnmount } from 'vue'

export function useTimeout() {
  let timeoutId = null

  const startTimeout = (callback, delay) => {
    clearTimeout(timeoutId)
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
