import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuthStore } from '../stores/useAuthStore';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requireSuperAdmin?: boolean;
  requireTenantAdmin?: boolean;
  requireAdminOrAbove?: boolean;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
  children,
  requireSuperAdmin,
  requireTenantAdmin,
  requireAdminOrAbove,
}) => {
  const { isAuthenticated, user } = useAuthStore();

  if (!isAuthenticated || !user) {
    return <Navigate to="/login" replace />;
  }

  const isSuperAdmin = Boolean(user.isSuperAdmin || user.roleName === 'super_admin' || user.role === 'super_admin');
  const isTenantAdmin = Boolean(user.isTenantAdmin || user.roleName === 'tenant_admin' || user.roleName?.includes('tenant_admin'));
  const isAdminOrAbove = isSuperAdmin || isTenantAdmin || user.role === 'admin';

  if (requireSuperAdmin && !isSuperAdmin) {
    return <Navigate to="/403" replace />;
  }

  if (requireTenantAdmin && !(isSuperAdmin || isTenantAdmin)) {
    return <Navigate to="/403" replace />;
  }

  if (requireAdminOrAbove && !isAdminOrAbove) {
    return <Navigate to="/403" replace />;
  }

  return <>{children}</>;
};
