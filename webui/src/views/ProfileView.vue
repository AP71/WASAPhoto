<script setup>
import Photo from '../components/Photo.vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import ErrorMsg from '../components/ErrorMsg.vue'
import { compile } from 'vue';
</script>

<script>
export default {
    components: {
        LoadingSpinner,
        ErrorMsg,
    },
    data: function () {
        return {
            errormsg: null,
            errormsgDialog: null,
			loading: false,
            id: null,
            username: null,
            photos: [],
            nextPhotosPageId: 0,
            followers: 0,
            following: 0,
            photoCounter: 0,
            followed: null,
            banned: null,
            dialog: false,
            uploadFile: false,
            newUsername: "",
        }
    },
    methods: {
        async getUserData() {
			this.loading = true;
			this.errormsg = null;
			try{

                const nextPage = this.nextPhotosPageId;
				let response = nextPage == 0?  
                                await this.$axios.get(`/profiles/${this.username}/`) : 
                                await this.$axios.get(`/profiles/${this.username}/?pageId=${this.nextPhotosPageId}`);
                if (nextPage == 0) {    
                    this.photos = response.data.photos;
                    this.id = response.data.id;
                    this.nextPhotosPageId = response.data.nextPhotosPageId;
                    this.followers = response.data.followers;
                    this.following = response.data.following;
                    this.photoCounter = response.data.photoCounter;
                } else {
                    this.photos.push(...response.data.photos);
                    this.nextPhotosPageId = response.data.nextPhotosPageId;
                }
			} catch(e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
        },
        async getStatus() {
            this.loading = true;
			this.errormsg = null;

            if (this.username != this.$profile.username) {
                // status
                try{
                    let response = await this.$axios.get(`/profiles/${this.username}/followed/${this.$profile.username}`);
                    this.followed = true;
                } catch(e) {
                    if (e.response.status === 404) {
                        this.followed = false;
                    }
                }

                try{
                    let response = await this.$axios.get(`/profiles/${this.username}/banned/${this.$profile.username}`);
                    this.banned = true;
                } catch(e) {
                    if (e.response.status === 404) {
                        this.banned = false;
                    }
                }
                
                if (!this.banned ) {
                    this.getUserData();
                }
            } else {
                this.getUserData();
            }
			
            this.loading = false;
        },
        async follow() {
            this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.put(`/profiles/${this.username}/followed/${this.$profile.username}`);
				this.followed = true;
				this.followers += 1;
			} catch(e) {
				this.errormsg = e.toString();
				if (e.response.status == 409) {
					this.followed = true;
				}
			}
			this.loading = false;
        },
        async unfollow() {
            this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.delete(`/profiles/${this.username}/followed/${this.$profile.username}`);
				this.followed = false;
				this.followers -= 1;
			} catch(e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
        },
        async ban() {
            this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.put(`/profiles/${this.username}/banned/${this.$profile.username}`);
				this.banned = true;
                this.reset();
			} catch(e) {
				this.errormsg = e.toString();
				if (e.response.status == 409) {
					this.banned = true;
				}
			}
			this.loading = false;
        },
        async unban() {
            this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.delete(`/profiles/${this.username}/banned/${this.$profile.username}`);
				this.banned = false;
                this.getUserData();
			} catch(e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
        },
        async changeUsername() {
            this.loading = true;
			this.errormsgDialog = null;
			try{
				let response = await this.$axios.put(`/profiles/${this.username}/username`, { newUsername: this.newUsername});
                this.username = response.data.username;
                this.$profile.username = response.data.username;
                this.dialog = false;
			} catch(e) {
				if (e.response.status == 400) {
					this.errormsgDialog = "Invalid username";
				}
                if (e.response.status == 409) {
					this.errormsgDialog = "Username already in use";
				}
			}
			this.loading = false;
        },
        async upload() {
            this.loading = true;
			this.errormsgDialog = null;
			try{
                const fileInput = document.getElementById('fileInput');
                const file = fileInput.files[0];
				let response = await this.$axios.post(`/profiles/${this.username}/photos/`, file);
                this.getUserData();
                this.dialog = false;
			} catch(e) {
				if (e.response.status == 400) {
					this.errormsgDialog = "Invalid file";
				}
                if (e.response.status == 413) {
					this.errormsgDialog = "File to large";
				}
                if (e.response.status == 415) {
					this.errormsgDialog = "File fromat not supported";
				}
			}
			this.loading = false;
        },
        reset() {
            this.photos = [];
            this.photoCounter = 0;
            this.followers = 0;
            this.following = 0;
        },
        showDialog(uploadFile) {
            this.uploadFile = uploadFile;
            this.dialog = true;
        },
        closeDialog() {
            this.dialog = false;
        },
    },
    mounted() {
        if (this.$profile.identifier==null){
			this.$router.push({path: '/'});
		} else {
            this.username = this.$route.params.username;
            this.getStatus();
		}
    },
}
</script>

<template>
    <div v-if="this.dialog" class="overlay"/>
    <div v-if="this.dialog && !this.uploadFile" class="dialogInterface" style="background-color: #383838;">
        <div class="fs-2 fw-bold pt-3 text-white">
            Change username
        </div>
        <div class="pt-2 pb-3">
            <input class="border border-3 border-success rounded-pill min-vh-25 min-vw-25 fs-4 text-indent" v-model="this.newUsername" placeholder="Username"/>
        </div>
        <ErrorMsg v-if="this.errormsgDialog" :msg="this.errormsgDialog"/>
        <div class="d-flex flex-row justify-content-evenly align-items-center pb-3 w-50">
            <button type="button" class="btn btn-outline-success rounded-pill fs-5" :disabled="this.newUsername.length<3 || this.newUsername.length>16" style="width: 150px" @click="this.changeUsername">
                Submit
            </button>
            <button type="button" class="btn btn-outline-danger rounded-pill fs-6" style="width: 150px" @click="this.closeDialog">
                Cancel
            </button>
        </div>
    </div>
    <div v-else-if="this.dialog && this.uploadFile" class="dialogInterface" style="background-color: #383838;">
        <div class="fs-2 fw-bold pt-3 text-white">
            Upload photo
        </div>
        <div class="pt-2 pb-3 text-white d-flex justify-content-center align-items-center">
            <input type="file" class="min-vh-25 min-vw-25 fs-4 text-indent" id="fileInput" />
        </div>
        <ErrorMsg v-if="this.errormsgDialog" :msg="this.errormsgDialog"/>
        <div class="d-flex flex-row justify-content-evenly align-items-center pb-3 w-50">
            <button type="button" class="btn btn-outline-success rounded-pill fs-5" style="width: 150px" @click="this.upload">
                Upload
            </button>
            <button type="button" class="btn btn-outline-danger rounded-pill fs-6" style="width: 150px" @click="this.closeDialog">
                Cancel
            </button>
        </div>
    </div>  
    <div class="d-flex justify-content-center align-items-start min-vh-100 w-100" style="background-color: #383838">
        <div class="d-flex flex-column align-items-center min-vh-100 w-75 text-white" style="background-color: #2e2e2e;">		
			<div style="height: 60px"/>
            <ErrorMsg v-if="errormsg" :msg="errormsg"/>
            <div class="fw-bolder text-white fs-1 d-flex justify-content-center align-items-center">
                {{ this.username }}
            </div>
            <div class="d-flex flex-row justify-content-evenly align-items-center w-100 fs-4 pt-4">
                <div class="d-flex flex-column justify-content-center align-items-center w-25">
                    <div class="fw-bold fs-3">
                        {{ this.photoCounter }}
                    </div>
                    Photo
                </div>
                <div class="d-flex flex-column justify-content-center align-items-center w-25">
                    <div class="fw-bold fs-3">
                        {{ this.followers}}
                    </div>
                    Followers
                </div>
                <div class="d-flex flex-column justify-content-center align-items-center w-25">
                    <div class="fw-bold fs-3">
                        {{ this.following }}
                    </div>
                    Followings
                </div>
            </div>
            <div class="d-flex flex-row justify-content-evenly align-items-center w-100 fs-4 pt-4">
                <div v-if="this.id!=this.$profile.identifier" class="d-flex flex-row justify-content-evenly align-items-center w-100">
                    <button v-if="!this.followed" @click="follow" type="button" class="btn btn-outline-success rounded-pill fs-4 w-25" style="width: 150px">
                        Segui
                    </button>
                    <button v-else type="button" @click="unfollow" class="btn btn-outline-danger rounded-pill fs-4 w-25" style="width: 150px">
                        Non seguire
                    </button>
                    <button v-if="!this.banned" @click="ban" type="button" class="btn btn-outline-danger rounded-pill fs-4 w-25" style="width: 150px">
                        Blocca
                    </button>
                    <button v-else type="button" @click="unban" class="btn btn-outline-success rounded-pill fs-4 w-25" style="width: 150px">
                        Sblocca
                    </button>
                </div>
                <div v-else class="d-flex flex-row justify-content-evenly align-items-center w-100">
                    <button type="button" @click="this.showDialog(false)" class="btn btn-outline-success rounded-pill fs-4 w-25" style="width: 150px">
                        Change username
                    </button>
                    <button type="button" @click="this.showDialog(true)" class="btn btn-outline-success rounded-pill fs-4 w-25" style="width: 150px">
                        Upload photo
                    </button>                        
                </div>
                <LoadingSpinner v-if="loading"/>
            </div>
            <div>
                
            </div>
		</div>
    </div>
</template>

<style scoped>
    .text-indent {
        text-indent: 20px;
    }

    .overlay {
        min-height: 100%;
        min-width: 100%;
        background-color: black;
        opacity: 70%;
        display: flex;
        justify-content: center;
        align-items: center;
        position: absolute;
        left: 50%;
        top: 50%;
        transform: translate(-50%, -50%);
        z-index: 9998;
    }

    .dialogInterface {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        border-radius: 2rem;
        min-width: 50%;
        opacity: 100%;
        position: absolute;
        left: 50%;
        top: 50%;
        transform: translate(-50%, -50%);
        z-index: 9999;
    }
</style>
