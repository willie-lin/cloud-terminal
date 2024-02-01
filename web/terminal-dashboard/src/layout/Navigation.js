import React from 'react';
import { Link } from 'react-router-dom';

const Navigation = ({ isLoggedIn, onLogout }) => {
    return (
        <div className="fixed w-full z-10">
            <div className="mx-auto px-6">
                <div className="flex justify-between items-center py-2">
                    <Link to="/">
                        {/*<Icon name="home" size="5xl" color="lightBlue"/>*/}
                        <span className="sr-only">Your Company</span>
                        <img className="h-8 w-auto" src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=600" alt=""/>
                    </Link>
                    {/*<div className="flex space-x-2">*/}
                    {/*    {isLoggedIn ? (*/}
                    {/*        <>*/}
                    {/*            <Button*/}
                    {/*                color="lightBlue"*/}
                    {/*                buttonType="filled"*/}
                    {/*                size="regular"*/}
                    {/*                rounded={false}*/}
                    {/*                block={false}*/}
                    {/*                iconOnly={false}*/}
                    {/*                ripple="light"*/}
                    {/*                onClick={onLogout}*/}
                    {/*            >*/}
                    {/*                Logout*/}
                    {/*            </Button>*/}
                    {/*        </>*/}
                    {/*    ) : (*/}
                    {/*        <>*/}
                    {/*            <Link to="/login">*/}
                    {/*                <Button*/}
                    {/*                    color="lightBlue"*/}
                    {/*                    buttonType="filled"*/}
                    {/*                    size="regular"*/}
                    {/*                    rounded={false}*/}
                    {/*                    block={false}*/}
                    {/*                    iconOnly={false}*/}
                    {/*                    ripple="light"*/}
                    {/*                >*/}
                    {/*                    Log in*/}
                    {/*                </Button>*/}
                    {/*            </Link>*/}
                    {/*            <Link to="/register">*/}
                    {/*                <Button*/}
                    {/*                    color="lightBlue"*/}
                    {/*                    buttonType="filled"*/}
                    {/*                    size="regular"*/}
                    {/*                    rounded={false}*/}
                    {/*                    block={false}*/}
                    {/*                    iconOnly={false}*/}
                    {/*                    ripple="light"*/}
                    {/*                >*/}
                    {/*                    Register*/}
                    {/*                </Button>*/}
                    {/*            </Link>*/}
                    {/*        </>*/}
                    {/*    )}*/}
                    {/*</div>*/}
                </div>
            </div>
        </div>
    );
};

export default Navigation;
