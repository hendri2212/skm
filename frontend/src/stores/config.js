import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useConfigStore = defineStore('config', () => {
    function detectApiBase() {
        const { hostname } = window.location
        if (hostname === 'localhost' || hostname === '127.0.0.1') {
            // Jika akses dari localhost/127.0.0.1
            return 'http://localhost:8080'
        } else {
            // Jika akses dari domain apa pun
            return 'https://api-skm.tukarjual.com'
        }
    }

    const apiHost = ref(`${detectApiBase()}/api`)

    function setApiHost(url) {
        apiHost.value = url
    }

    return {
        apiHost,
        setApiHost
    }
})