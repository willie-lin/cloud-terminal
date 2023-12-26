import LoginForm from '../components/LoginForm';

export default function Login({ onLogin }) {
    return (
        <div>
            <LoginForm onLogin={onLogin} />
        </div>
    );
}
