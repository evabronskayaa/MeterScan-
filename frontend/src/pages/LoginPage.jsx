import { Anchor, Button, Container, Input, Paper, PasswordInput, PinInput, Text, TextInput, Title } from '@mantine/core'
import { Form, useForm } from '@mantine/form'
import { useState } from 'react'
import { Link } from 'react-router-dom'
import useApp from '../hooks/useApp.tsx'
import useInvisibleRecaptcha from '../hooks/useInvisibleRecaptcha.tsx'
import { useTitle } from '../hooks/useTitle.tsx'
import { checkResponse, fetchApi, setToken } from '../lib/api.ts'
import Notifications from '../lib/notifications.tsx'

const LoginPage = () => {
    const { updateApp } = useApp()
    const [loading, setLoading] = useState(false)
    const form = useForm({
        validateInputOnChange: true,
        initialValues: {
            username: '', password: '', code: ''
        },

        validate: {
            username: value => {
                if (value.length < 3) {
                    return 'Имя пользователя должно быть больше 2 символов'
                }
                if (value.length > 20) {
                    return 'Имя пользователя должно быть не больше 20 символов'
                }
                return null
            },
            password: value => value.length < 6 ? 'Пароль должен быть больше 5 символов' : null
        }
    })
    const { recaptchaComponent, getRecaptchaValue } = useInvisibleRecaptcha()
    const [needCode, setNeedCode] = useState(false)

    const onSubmit = async (values) => {
        if (loading) {
            return
        }
        setLoading(true)
        getRecaptchaValue().then(recaptcha => {
            fetchApi('/sessions', {
                method: 'POST', body: {
                    username: values.username,
                    password: values.password,
                    code: values.code,
                    recaptcha: recaptcha
                }
            }).then(checkResponse)
                .then(r => r.json())
                .then(data => {
                    if (data.need_code) {
                        setNeedCode(true)
                    } else if (data.token) {
                        Notifications.success({ message: 'Вход успешно выполнен' })
                        setToken(data.token)
                        updateApp({
                            user: data.user
                        })
                    }
                })
                .catch(error => {
                    Notifications.error({ message: error })
                })
                .finally(() => setLoading(false))
        })
    }

    useTitle('Авторизация')

    const disabledForm = needCode || loading

    return <Container size={ 420 } py={ 40 }>
            <Title ta="center" className="title">
                С возвращением!
            </Title>
            <Text c="dimmed" size="sm" ta="center" mt={ 5 }>
                Еще нет аккаунта? <Anchor component={ Link } to="/register">Создать аккаунт</Anchor>
            </Text>

        <Paper withBorder shadow="md" p={ 30 } mt={ 30 }>
                <Form form={ form } onSubmit={ onSubmit }>
                    <TextInput label="Имя пользователя" required disabled={ disabledForm }
                               { ...form.getInputProps('username') }/>
                    <PasswordInput label="Пароль" required mt="md" disabled={ disabledForm }
                                   { ...form.getInputProps('password') }/>
                    { needCode &&
                        <Input.Wrapper mt="sm" label="Код из приложения">
                            <PinInput type="number" onComplete={ () => onSubmit(form.values) }
                                      length={ 6 } autoFocus { ...form.getInputProps('code') }/>
                        </Input.Wrapper>
                    }

                    { recaptchaComponent }

                    <Button type="submit" fullWidth mt="xl" loading={ loading }
                            disabled={ !form.isValid() || needCode }>
                        Войти
                    </Button>
                </Form>
            </Paper>
    </Container>
}

export default LoginPage