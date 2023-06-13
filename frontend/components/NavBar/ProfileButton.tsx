import React, { MouseEventHandler } from 'react';
import {
    Avatar,
    Button,
    Menu,
    MenuButton,
    MenuList,
    MenuItem,
    MenuDivider,
    useToast
} from '@chakra-ui/react';
import axios from 'axios';

interface ProfileButtonProps {
    profilePic: string;
}

export default function ProfileButton(props: ProfileButtonProps) {
    const {
        profilePic = 'https://images.unsplash.com/photo-1493666438817-866a91353ca9?ixlib=rb-0.3.5&q=80&fm=jpg&crop=faces&fit=crop&h=200&w=200&s=b616b2c5b373a80ffc9636ba24f7a4a9'
    } = props;
    const toast = useToast();

    const handleClick: MouseEventHandler = () => {
        var url = 'http://localhost:8080/auth/logout'
        axios.get(url, {withCredentials: true})
            .then((res) => {
                console.log(res);
                toast({
                    title: "Form submission successful.",
                    description: "We've received your form data.",
                    status: "success",
                    duration: 5000,
                    isClosable: true,
                });
            }).catch((error) => {
                console.log(error.response);
                toast({
                    title: "An error occurred.",
                    description: error.response.data.error,
                    status: "error",
                    duration: 5000,
                    isClosable: true,
                });
            });
    };

    return (
        <Menu>
                            
            <MenuButton
                as={Button}
                rounded={'full'}
                variant={'link'}
                cursor={'pointer'}
                minW={0}
            >
                <Avatar
                size={'md'}
                src={profilePic}
                />
            </MenuButton>
            <MenuList>
                <MenuItem>View your profile</MenuItem>
                <MenuDivider />
                <MenuItem><Button onClick={handleClick}>Log out</Button></MenuItem>
            </MenuList>
        </Menu>
    );
}