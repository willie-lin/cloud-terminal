import RegisterForm from '../components/user/RegisterForm';

export default function Register({ onRegister }) {
    return (
        <div>
            <RegisterForm onRegister={onRegister} />
        </div>
    );
}
