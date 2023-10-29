import { MantineColor } from '@mantine/core'
import { showNotification } from '@mantine/notifications'
import { IconCheck, IconX } from '@tabler/icons-react'
import React from 'react'

let idCounter = 0
const genNextId = () => {
    idCounter++
    return '' + idCounter
}

interface NotificationOptions {
    title?: React.ReactNode
    message: React.ReactNode
    withCloseButton?: boolean
    autoClose?: number
    color?: MantineColor
    icon?: React.ReactNode
}

const Notifications = {
    add(options: NotificationOptions) {
        return showNotification({
            id: genNextId(),
            autoClose: 5000,
            ...options
        })
    },
    error(options: NotificationOptions) {
        return this.add({
            color: 'red',
            icon: <IconX/>,
            ...options
        })
    },
    success(options: NotificationOptions) {
        return this.add({
            color: 'green',
            icon: <IconCheck/>,
            ...options
        })
    }
}

export default Notifications