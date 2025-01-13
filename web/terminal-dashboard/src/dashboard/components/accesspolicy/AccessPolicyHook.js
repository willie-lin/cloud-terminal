import {useEffect, useState} from "react";
import { getAllAccessPolicies } from "../../../api/api";

export const useFetchAccessPolicies = () => {
    const [accessPolicies, setAccessPolicies] = useState([]);

    useEffect(() => {
        getAllAccessPolicies()
            .then(data => setAccessPolicies(data))
            .catch(error => console.error('Error:', error));
    }, []);
    return accessPolicies;
};
