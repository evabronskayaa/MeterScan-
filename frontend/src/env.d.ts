interface ImportMetaEnv {
    readonly VITE_API_ENDPOINT: string
    readonly VITE_RECAPTCHA_KEY: string
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}

declare module '*.svg' {
    const content: React.FunctionComponent<React.SVGAttributes<SVGElement>>
    export default content
}