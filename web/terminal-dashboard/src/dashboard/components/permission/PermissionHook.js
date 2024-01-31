import {useEffect, useState} from "react";
import {getAllPermissions} from "../../../api/api";


export const useFetchPermissions = () => {
    const [permissions, setPermissions] = useState([]);

    useEffect(() => {
        getAllPermissions()
            .then(data => setPermissions(data))
            .catch(error => console.error('Error:', error));
    }, []);
    return permissions;
};
