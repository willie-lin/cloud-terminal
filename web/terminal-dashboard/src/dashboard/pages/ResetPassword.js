import ResetPasswordForm from "../components/user/ResetPasswordForm";

export default function ResetPassword({ onResetPassword }) {
    return (
        <div>
            <ResetPasswordForm onResetPassword={onResetPassword} />
        </div>
    );
}