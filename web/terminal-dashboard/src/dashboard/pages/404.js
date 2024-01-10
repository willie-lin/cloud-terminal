import React from 'react';
import {Button} from "@material-tailwind/react";
import {Link} from "react-router-dom";

function NotFoundPage() {
    return (
        <div className="flex flex-col items-center justify-center h-screen text-center">
            <h1 className="mt-4 text-6xl font-semibold">
                404 😢
            </h1>
            <p className="mt-4 text-lg">
                Lost in the Digital Wilderness
            </p>
            <p className="mt-4 text-md">
                You’ve ventured into uncharted digital territory. The page you seek has eluded us. Let’s guide you back to familiar paths.
            </p>
            <Link to="/">
                <Button
                color="lightBlue"
                buttonType="filled"
                size="lg"
                rounded={false}
                block={false}
                iconOnly={false}
                ripple="light"
                className="mt-4"
                >
                    BACK TO HOME
                </Button>
            </Link>
        </div>
    );
}

export default NotFoundPage;
