/**
 * 全局消息通知服务
 * 提供统一的消息提示功能，可在应用的任何地方调用
 */
import { createApp, h } from 'vue'
import AppMessage from '../components/ui/AppMessage.vue'

// 默认配置
const defaultOptions = {
  type: 'info',
  duration: 3000,
  closable: true,
  showIcon: true,
  autoClose: true,
  filled: false
}

// 创建消息容器
const createMessageContainer = () => {
  const container = document.createElement('div')
  container.className = 'app-message-container fixed top-4 right-4 z-50 flex flex-col gap-2 w-80'
  document.body.appendChild(container)
  return container
}

// 获取或创建消息容器
const getMessageContainer = () => {
  let container = document.querySelector('.app-message-container')
  if (!container) {
    container = createMessageContainer()
  }
  return container
}

// 显示消息
const showMessage = (options) => {
  const container = getMessageContainer()
  const messageNode = document.createElement('div')
  container.appendChild(messageNode)

  // 合并选项
  const messageOptions = {
    ...defaultOptions,
    ...options
  }

  // 创建消息实例
  const app = createApp({
    render() {
      return h(AppMessage, {
        ...messageOptions,
        onClose: () => {
          // 移除消息
          setTimeout(() => {
            app.unmount()
            container.removeChild(messageNode)
            if (container.children.length === 0) {
              document.body.removeChild(container)
            }
          }, 300) // 延迟移除，以便过渡动画完成
        }
      }, () => messageOptions.content)
    }
  })

  app.mount(messageNode)

  return {
    close: () => {
      app.unmount()
      if (messageNode.parentNode) {
        container.removeChild(messageNode)
        if (container.children.length === 0) {
          document.body.removeChild(container)
        }
      }
    }
  }
}

// 导出不同类型的消息方法
export default {
  info(content, options = {}) {
    return showMessage({ type: 'info', content, ...options })
  },
  success(content, options = {}) {
    return showMessage({ type: 'success', content, ...options })
  },
  warning(content, options = {}) {
    return showMessage({ type: 'warning', content, ...options })
  },
  error(content, options = {}) {
    return showMessage({ type: 'error', content, ...options })
  },
  // 自定义消息
  custom(options) {
    return showMessage(options)
  }
}