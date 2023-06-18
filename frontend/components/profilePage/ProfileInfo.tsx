import {
    Text,
    Box,
    VStack,
} from '@chakra-ui/react';
import InfoSection from './InfoSection';
import ProjectDisplay from './ProjectDisplay';
import { useEffect, useState } from 'react';
import axios from 'axios';
import { User } from '../../types';
import { useUser } from '../../userContext';


export default function ProfileInfo({username}: {username: string}) {
   
    const [user, setUser] = useState<User>({
        AboutMe: "",
        Email: "",
        Name: "",
        Title: "",
        ProfilePic: "",
        Username: "",
        Projects: []
    });
    const [profileState, setProfileState] = useState<string>("loading");

    useEffect(() => {
        const currentUrl = "http://localhost:8080/auth/user";
        const profileUrl = "http://localhost:8080/users/" + username;
        if (username) {
            axios.get(currentUrl, { withCredentials: true })
            .then((res) => {
                const currentUser = res.data.data.Username;
                axios.get(profileUrl).then((res) => {
                    const { AboutMe, Email, Name, Title, ProfilePic, Projects } = res.data.data;
                    setUser({
                        AboutMe: AboutMe ? AboutMe : "No description available",
                        Email: Email,
                        Name: Name ? Name : "No display name",
                        Title: Title ? Title : "No title available",
                        ProfilePic: ProfilePic,
                        Username: username,
                        Projects: Projects,
                    });
                    // Compare profile user to current user
                    
                    if (currentUser === username) {
                        setProfileState("self");
                    } else {
                        setProfileState("other");
                    }
                }).catch((err) => {    
                    setProfileState("invalid")
                });
            })
            .catch((err) => {
                console.log(err);
            })
        }
        // Fetch current user
    }, [username]); 

    if (profileState === "invalid") {
        return (
            <p>invalid username</p>
        )
    }

    return (
        <Box mt={10} mx={5} p={4}>
          <VStack spacing={10} align="start">
            <InfoSection
              user={user}
              {...(profileState === "self" && { setUser })}
            />
            <ProjectDisplay />
          </VStack>
        </Box>
      );
};
