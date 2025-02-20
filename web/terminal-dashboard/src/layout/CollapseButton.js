import React from 'react';
import { IconButton } from "@material-tailwind/react";
import { ChevronLeftIcon, ChevronRightIcon } from 'lucide-react';

const CollapseButton = ({ isCollapsed, onCollapse, isDarkMode }) => {
    return (
        <IconButton
            variant="text"
            onClick={onCollapse}
        >
            {isCollapsed ?
                <ChevronRightIcon className={`h-5 w-5 ${isDarkMode ? 'text-gray-500' : 'text-blue-gray-500'}`} /> :
                <ChevronLeftIcon className={`h-5 w-5 ${isDarkMode ? 'text-gray-500' : 'text-blue-gray-500'}`} />
            }
        </IconButton>
    );
};

export default CollapseButton;
