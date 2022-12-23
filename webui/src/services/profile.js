
const Profile = {
    identifier: null,
    username: null,
    isLogged: false,

    setProfile(identifier, username) {
        this.identifier = identifier;
        this.username = username;
        this.isLogged = true;
    },

    logout() {
        this.identifier = null;
        this.username = null;
        this.isLogged = false;
    }
    
}

export default Profile;