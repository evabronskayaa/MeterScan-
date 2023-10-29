import { Anchor, Button, Container, Paper, PasswordInput, Text, TextInput, Title } from '@mantine/core'
import { Form, useForm } from '@mantine/form'
import { useState } from 'react'
import { Link } from 'react-router-dom'
import useApp from '../hooks/useApp'
import useInvisibleRecaptcha from '../hooks/useInvisibleRecaptcha'
import { useTitle } from '../hooks/useTitle'
import { checkResponse, fetchApi, setToken } from '../lib/api'
import Notifications from '../lib/notifications'

const RegisterPage = () => {
    const form = useForm({
        initialValues: {
            username: '', password: '', confirmPassword: '', key: ''
        }, validateInputOnChange: true, validate: {
            username: value => {
                if (value.length < 3) {
                    return 'Имя пользователя должно быть больше 2 символов'
                }
                if (value.length > 20) {
                    return 'Имя пользователя должно быть не больше 20 символов'
                }
                return null
            },
            password: value => value.length < 6 ? 'Пароль должен быть больше 5 символов' : null,
            confirmPassword: (value, values) => {
                if (value !== values.password) {
                    return 'Пароль не совпадает'
                }
                return null
            }
        }
    })
    const { recaptchaComponent, getRecaptchaValue } = useInvisibleRecaptcha()
    const { updateApp } = useApp()
    const [loading, setLoading] = useState(false)

    const onSubmit = async (values: any) => {
        if (loading) {
            return
        }
        setLoading(true)
        getRecaptchaValue().then(recaptchaValue => {
            fetchApi('/users', {
                method: 'POST', body: {
                    username: values.username, password: values.password, recaptcha: recaptchaValue, key: values.key
                }
            }).then(checkResponse)
                .then(r => r.json())
                .then(data => {
                    if (data.token) {
                        Notifications.success({ message: 'Регистрация успешно выполнена' })
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

    useTitle('Регистрация')

    return <Container size={ 420 } py={ 40 }>
        <Title ta="center" className="title">
            Добро пожаловать!
        </Title>
        <Text c="dimmed" size="sm" ta="center" mt={ 5 }>
            Уже есть аккаунт? <Anchor component={ Link } to="/login">Войти</Anchor>
        </Text>

        <Paper withBorder shadow="md" p={ 30 } mt={ 30 }>
            <Form form={ form } onSubmit={ onSubmit }>
                <TextInput label="Имя пользователя" required
                           { ...form.getInputProps('username') }/>
                <PasswordInput label="Пароль" required mt="md"
                               { ...form.getInputProps('password') }/>
                <PasswordInput label="Повторить пароль" required mt="md"
                               { ...form.getInputProps('confirmPassword') }/>
                <TextInput label="Ключ" required mt="md"
                           { ...form.getInputProps('key') }/>

                { recaptchaComponent }

                <Button type="submit" fullWidth mt="xl" loading={ loading }>
                    Зарегистрироваться
                </Button>
            </Form>
        </Paper>
    </Container>
}

export default RegisterPage