<script setup>
import Photo from '../components/Photo.vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import ErrorMsg from '../components/ErrorMsg.vue'
</script>

<script>
export default {
    components: {
        LoadingSpinner,
        ErrorMsg,
    },
    data: function () {
        return {
            userToSearch: null,
            users: [],
            nextPageId: 0,
            errormsg: null,
			loading: false,
        }
    },
    methods: {
        async getUsers() {
            this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.get(`/profiles/?userToSearch=${this.userToSearch}`);
				this.users.push(...response.data.users);
				this.nextPageId = response.data.nextUsersPageId;
                this.userToSearch = null;
			} catch(e) {
				this.errormsg = e.toString();
			}
            
			this.loading = false;
        },
        openProfile(user) {
            this.$router.push({ name: `profiles`, params: {username: user.username} });
        }

    },
    mounted() {
        if (this.$profile.identifier == null) {
            this.$router.push({ path: '/' });
        } else {

        }
    }
}
</script>

<template>
    <div class="d-flex justify-content-center align-items-start min-vh-100 w-100" style="background-color: #383838">
        <div class="d-flex flex-column justify-content-centeralign-items-start min-vh-100 w-75" style="background-color: #2e2e2e;">		
			<div class="cerca">
                <input v-model="this.userToSearch" class="inputBar border border-2 border-success fs-4" style="background-color: #383838;" placeholder="Username"/>
                <LoadingSpinner :loading="this.loading" />
                <i class="ps-4 icon bi bi-send-fill text-success fs-4" @click="getUsers"></i>
            </div>
            <div class="d-flex flex-column justify-content-center align-items-center pt-4">
                <div v-for="user in this.users" key="user.identifier" class="result" @click="openProfile(user)">
                    <div class="d-flex justify-content-center align-items-start min-vw-100">
                        <div class="resultMod fw-bold fs-4">
                            {{ user.username }}
                        </div>
                    </div>
                    <div style="height: 50px; background-color: #2e2e2e;"/>
                </div>
            </div>
		</div>
    </div>
</template>

<style scoped>

    .result {
        display: flex;
        justify-content: center;
        align-items: center;
        min-width: 50%;
    }

    .resultMod {
        display: flex;
        justify-content: start;
        text-indent: 20px;
        align-items: center;
        color: white;
        background-color: #212121;
        font-size: 20px;
        min-width: 45%;
        height: 40px;
        border-radius: 2rem;
    }

    .cerca {
        display: flex;
        justify-content: center;
        align-items: center;
        min-width: 100%;
        padding-top: 75px;
    }

    .inputBar {
        border-radius: 2rem;
        min-width: 60%;
        color: white;
        text-indent: 10px;
    }
</style>
