import React, { useState } from 'react';
import { useDisclosure, useToast } from '@chakra-ui/react';
import EditPicButton from './EditPicButton';
import ImageCropper from './ImageCropper';
import axios from 'axios';
import { User } from '../../../types';
import { useUser } from '../../../userContext';


interface CropperComponentProps {
    profilePic?: string;
    setUser: React.Dispatch<React.SetStateAction<User>>;
}

const CropperComponent: React.FC<CropperComponentProps> = ({profilePic, setUser}) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const [selectedImage, setSelectedImage] = useState<string | undefined>(undefined);
    const toast = useToast();
    const { setNeedUpdate } = useUser();
  
    const handleValidFile = (file: File) => {
        setSelectedImage(URL.createObjectURL(file));
        onOpen();
    };
  
    const handleCroppedImage = (dataUrl: string) => {
        const url = "http://localhost:8080/auth/user"
        axios.patch(url, {
                profilepic: dataUrl,
            }, {
                withCredentials: true,
            })
            .then(res => {
                const { ProfilePic } = res.data.data;
                setUser((prevUser: User) => ({
                    ...prevUser,
                    ProfilePic: ProfilePic
                }));
                setNeedUpdate(true);
                toast({
                    title: "Profile updated.",
                    description: "We've updated your profile picture for you.",
                    status: "success",
                    duration: 9000,
                    isClosable: true,
                });
            })
            .catch(error => {
                console.error(error);
            });
    };

    return (
        <>
            <EditPicButton currentProfilePic={profilePic || ""} onValidFile={handleValidFile} />
            <ImageCropper isOpen={isOpen} onClose={onClose} onCropped={handleCroppedImage} imageSrc={selectedImage} />
        </>
    );
};

export default CropperComponent;
