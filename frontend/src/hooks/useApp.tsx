import { MantineTheme } from '@mantine/core'
import { createContext, useContext } from 'react'
import { PartialDeep } from 'type-fest'
import { Agent } from '../lib/models'

interface AppContextContainerType {
    app: AppContextType,
    updateApp: (changes: AppContextTypeOverride) => void
}

interface AppContextType {
    user?: Agent,
    theme: MantineTheme
}

export type AppContextTypeOverride = PartialDeep<AppContextType>

export const AppContext = createContext<AppContextContainerType>(null!)

export default () => useContext(AppContext)