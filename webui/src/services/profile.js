
const Profile = {
    identifier: null,
    username: null,
    

    setProfile(identifier, username) {
        this.identifier = identifier;
        this.username = username;
    },

    isLogged() {
        return this.identifier!=null && this.username!=null;
    }
    
}

export default Profile;