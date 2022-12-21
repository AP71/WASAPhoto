
const Profile = {
    identifier: null,
    username: null,
    isLogged: false,

    setProfile(identifier, username) {
        this.identifier = identifier;
        this.username = username;
        this.isLogged = true;
    },

    
}

export default Profile;