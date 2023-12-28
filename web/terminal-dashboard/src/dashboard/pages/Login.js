import LoginForm from '../components/user/LoginForm';

export default function Login({ onLogin }) {
    return (
        <div>
            <LoginForm onLogin={onLogin} />
        </div>
    );
}
