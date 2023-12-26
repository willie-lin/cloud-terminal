import RegisterForm from '../components/RegisterForm';

export default function Register({ onRegister }) {
    return (
        <div>
            <RegisterForm onRegister={onRegister} />
        </div>
    );
}
