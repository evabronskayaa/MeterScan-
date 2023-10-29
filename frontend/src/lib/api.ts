import { TokenInfo } from './models'

let cachedToken: TokenInfo | undefined

interface ApiRequestInit extends Omit<RequestInit, 'body' | 'headers'> {
    body?: BodyInit | object | null
    headers?: Record<string, string>
}

function setHeader(options: ApiRequestInit, name: string, value: string) {
    options.headers = options.headers || {}
    options.headers[name] = value
}

// https://stackoverflow.com/a/8511350/6620659
function isObject(val: unknown): boolean {
    return typeof val === 'object' &&
        !Array.isArray(val) &&
        val !== null
}

export const fetchApi = async (path: string, options: ApiRequestInit = {}) => {
    if (cachedToken) {
        setHeader(options, 'Authorization', 'Bearer ' + cachedToken.token)
    }

    if (isObject(options.body) && !(options.body instanceof FormData)) {
        options.body = JSON.stringify(options.body)
        setHeader(options, 'Content-Type', 'application/json')
    }

    setHeader(options, 'Accept', 'application/json')

    const url = import.meta.env.VITE_API_ENDPOINT + path
    return fetch(url, options as RequestInit)
}

export const checkResponse = async (r: Response) => {
    if (r.ok) {
        return r
    }
    let body
    try {
        body = await r.json()
    } catch {
        throw 'Невозможно подключиться к серверу'
    }
    if (body.error) {
        throw body.error
    } else {
        throw 'Некорректный ответ сервера'
    }
}

export const setToken = (token?: TokenInfo) => {
    cachedToken = token
    if (token)
        localStorage.setItem('at', JSON.stringify(token))
    else
        localStorage.removeItem('at')
}

export const getToken = (): TokenInfo | undefined => {
    if (cachedToken == null) {
        const tokenFromStorage = localStorage.getItem('at')
        if (tokenFromStorage != null) {
            try {
                cachedToken = JSON.parse(tokenFromStorage) as TokenInfo
            } catch {
                setToken(undefined)
            }
        }
    }
    return cachedToken
}

export const getBase64Image = async (res: Response) => {
    const blob = await res.blob()

    const reader = new FileReader()

    await new Promise((resolve, reject) => {
        reader.onload = resolve
        reader.onerror = reject
        reader.readAsDataURL(blob)
    })
    return reader.result + ''
}