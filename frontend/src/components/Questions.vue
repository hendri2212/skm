<template>
    <div class="question-wrapper mx-auto" style="max-width: 414px;">
        <!-- ===== FORM DATA PENGGUNA ===== -->
        <div v-if="!quizStarted" id="user-form">
            <div class="d-flex justify-content-between align-items-center" style="margin-bottom: 1.5rem;">
                <button class="btn btn-link p-0 bg-white rounded-circle" style="width: 48px; height: 48px;"
                    @click="backToHome()">
                    <i class="bi bi-arrow-left fs-4 text-secondary"></i>
                </button>
                <h6 class="mb-0">Data Pengguna</h6>
                <button class="btn p-0"></button>
            </div>
            <div class="card question-card mb-4 rounded-4 shadow-sm">
                <div class="card-body">
                    <!-- Nama Lengkap -->
                    <div class="form-floating mb-3">
                        <input v-model="user.name" type="text" class="form-control" id="input-name"
                            placeholder="Nama Lengkap" required />
                        <label for="input-name">Nama Lengkap</label>
                    </div>

                    <div class="row">
                        <!-- Tempat Lahir -->
                        <div class="col-md-6 mb-3">
                            <div class="form-floating">
                                <input v-model="user.birthPlace" type="text" class="form-control" id="input-birth-place"
                                    placeholder="Tempat Lahir" required />
                                <label for="input-birth-place">Tempat Lahir</label>
                            </div>
                        </div>

                        <!-- Tanggal Lahir -->
                        <div class="col-md-6 mb-3">
                            <div class="form-floating">
                                <input v-model="user.birthDate" type="date" class="form-control" id="input-birth-date"
                                    placeholder="Tanggal Lahir" required />
                                <label for="input-birth-date">Tanggal Lahir</label>
                            </div>
                        </div>
                    </div>

                    <!-- Jenis Kelamin (Radio Inputs) -->
                    <div class="mb-3">
                        <label class="form-label d-block">Jenis Kelamin</label>
                        <div class="form-check form-check-inline">
                            <input class="form-check-input" type="radio" name="is_male" id="gender-male"
                                :value="true" v-model="user.is_male" required />
                            <label class="form-check-label" for="gender-male">Laki-laki</label>
                        </div>
                        <div class="form-check form-check-inline">
                            <input class="form-check-input" type="radio" name="is_male" id="gender-female"
                                :value="false" v-model="user.is_male" required />
                            <label class="form-check-label" for="gender-female">Perempuan</label>
                        </div>
                    </div>

                    <!-- Pendidikan Terakhir -->
                    <div class="form-floating mb-3">
                        <select v-model="user.education" class="form-select" id="select-education" required>
                            <option value="" disabled>— Pilih Pendidikan Terakhir —</option>
                            <option v-for="edu in educations" :key="edu.id" :value="edu.id">
                                {{ edu.name }}
                            </option>
                        </select>
                        <label for="select-education">Pendidikan Terakhir</label>
                    </div>

                    <!-- Pekerjaan Utama -->
                    <div class="form-floating mb-3">
                        <select v-model="user.occupation" class="form-select" id="select-occupation" required>
                            <option value="" disabled>— Pilih Pekerjaan Utama —</option>
                            <option v-for="ecc in occupations" :key="ecc.id" :value="ecc.id">
                                {{ ecc.name }}
                            </option>
                        </select>
                        <label for="select-occupation">Pekerjaan Utama</label>
                    </div>

                    <!-- Tombol Mulai Quiz -->
                    <button class="btn btn-warning w-100 mt-5" :disabled="!canStartQuiz" @click="startQuiz">
                        Mulai Quiz
                    </button>
                </div>
            </div>
        </div>

        <!-- ===== LAYAR SOAL ===== -->
        <div v-else-if="!showResult" class="quiz-screen">
            <!-- Header: Nomor Soal -->
            <div class="d-flex justify-content-between align-items-center" style="margin-bottom: 1.5rem;">
                <button class="btn btn-link p-0 bg-white rounded-circle" @click="prevQuestion" style="width: 48px; height: 48px;">
                    <i class="bi bi-arrow-left fs-4 text-secondary"></i>
                </button>
                <h6 class="mb-0">
                    Soal {{ currentIndex + 1 }} / {{ questions.length }}
                </h6>
                <button class="btn p-0"></button>
            </div>

            <!-- Kartu Soal Utama -->
            <div class="card question-card mb-4 shadow-sm rounded-4">
                <div class="card-body">
                    <!-- Teks Soal -->
                    <p class="question-text p-5 rounded-4 text-center">{{ currentQuestion.question }}</p>

                    <!-- Progress Bar Manual -->
                    <div class="mb-3">
                        <div class="time-bar-container">
                            <div class="time-bar-fill"
                                :style="{ width: ((currentIndex + 1) / questions.length * 100) + '%' }"></div>
                        </div>
                    </div>

                    <!-- Pilihan Jawaban (A, B, C, D) -->
                    <div class="answer-list">
                        <div class="answer-option" v-for="(label, key) in currentQuestion.options" :key="key" :class="{
                            selected: selectedOption === key
                        }" @click="selectOption(key)">
                            <div class="d-flex align-items-center w-100">
                                <!-- Huruf Pilihan -->
                                <div class="option-label me-3">{{ key.toUpperCase() }}</div>
                                <!-- Teks Pilihan -->
                                <div class="option-text">{{ label }}</div>
                            </div>
                        </div>
                    </div>

                    <button class="btn btn-warning w-100 mt-5" @click="goToNext" :disabled="!selectedOption">
                        {{ isLastQuestion ? 'Lihat Hasil' : 'Soal Berikutnya' }}
                    </button>
                </div>
            </div>

        </div>

        <!-- ===== HALAMAN HASIL AKHIR ===== -->
        <div v-else class="result-screen text-center">
            <!-- Header: Nomor Soal -->
            <div class="d-flex justify-content-between align-items-center" style="margin-bottom: 1.5rem;">
                <button class="btn btn-link p-0 bg-white rounded-circle" @click="backToHome"
                    :disabled="currentIndex === 0" style="width: 48px; height: 48px;">
                    <i class="bi bi-house fs-4 text-secondary"></i>
                </button>
            </div>
            <div class="position-relative" style="margin-top: 60px;">
                <div class="card bg-white bg-opacity-75 rounded-4 position-absolute" style="width: 100%; height: 100%; top: -30px;">
                    <label class="py-1 fw-bold text-secondary">Survey Selesai</label>
                </div>
                <div class="card bg-white bg-opacity-90 rounded-4 shadow-sm p-4 position-relative">
                    <div class="text-center">
                        <img src="https://upload.wikimedia.org/wikipedia/commons/0/0a/Lambang_Kabupaten_Kotabaru.png" alt="User Image" class="mx-auto" style="width: 80px">
                        <h2 class="fs-4 fw-bold mt-2 text-secondary">Pulaulaut Sigam</h2>
                    </div>
                    <p class="fs-5 fw-bold text-dark my-4 text-center">Terima kasih atas partisipasi Anda</p>
                    <p class="small text-muted mb-4 text-center">Anda telah menyelesaikan quiz ini!</p>
                    <button class="btn btn-warning w-100 rounded-3 py-2" @click="restartQuiz">Ulangi Quiz</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { useConfigStore } from '@/stores/config'
export default {
    name: 'Questions',
    setup() {
        const config = useConfigStore()
        return { apiHost: config.apiHost }
    },
    data() {
        return {
            user: {
                name: '',
                birthPlace: '',
                birthDate: '',
                is_male: null,
                education: '',
                occupation: ''
            },
            quizStarted: false,
            questions: [],
            educations: [],
            occupations: [],
            selectedOption: '',
            currentIndex: 0,
            showResult: false,
            // mapping question.id -> choice_id
            userAnswers: {}
        };
    },
    computed: {
        currentQuestion() {
            return this.questions[this.currentIndex] || { options: {} };
        },
        isLastQuestion() {
            return this.currentIndex === this.questions.length - 1;
        },
        canStartQuiz() {
            return (
                this.user.name.trim() &&
                this.user.birthPlace.trim() &&
                this.user.birthDate &&
                this.user.is_male !== null &&
                this.user.education &&
                this.user.occupation
            );
        }
    },
    mounted() {
        this.fetchQuestions();
        this.fetchEducations();
        this.fetchOccupations();
    },
    methods: {
        backToHome() {
            this.$router.push({ name: 'home' });
        },
        startQuiz() {
            this.quizStarted = true;
            this.userAnswers = {};
        },
        selectOption(key) {
            this.selectedOption = key;
        },
        prevQuestion() {
            if (this.currentIndex === 0) {
                this.quizStarted = false;
            } else {
                this.currentIndex--;
                const prevId = this.questions[this.currentIndex].id;
                // retrieve previously selected letter
                this.selectedOption = Object.keys(this.currentQuestion.options)
                    .find(k => this.userAnswers[prevId] === this.currentQuestion.choiceIds[k]) || '';
            }
        },
        async goToNext() {
            const q = this.currentQuestion;
            // simpan choice_id berdasarkan letter
            this.userAnswers[q.id] = q.choiceIds[this.selectedOption];
            this.selectedOption = '';

            if (!this.isLastQuestion) {
                this.currentIndex++;
            } else {
                this.showResult = true;
                await this.submitAnswers();
            }
        },
        restartQuiz() {
            this.quizStarted = false;
            this.showResult = false;
            this.currentIndex = 0;
            this.selectedOption = '';
            this.userAnswers = {};
        },
        async fetchEducations() {
            try {
                const response = await fetch(`${this.apiHost}/educations`);
                const data = await response.json();
                this.educations = data.educations;
            } catch (err) {
                console.error('Gagal memuat data pendidikan:', err);
                this.error = 'Gagal memuat data pendidikan.';
            } finally {
                this.loading = false;
            }
        },
        async fetchOccupations() {
            try {
                const response = await fetch(`${this.apiHost}/occupations`);
                const data = await response.json();
                this.occupations = data.occupations;
            } catch (err) {
                console.error('Gagal memuat data pekerjaan:', err);
                this.error = 'Gagal memuat data pekerjaan.';
            } finally {
                this.loading = false;
            }
        },
        async fetchQuestions() {
            try {
                const res = await fetch(`${this.apiHost}/questions`)
                const data = await res.json();
                this.questions = data.map(q => {
                    const sorted = q.choices.sort((a, b) => a.points - b.points);
                    return {
                        id: q.id,
                        question: q.question_text,
                        options: {
                            a: sorted[0].choice_text,
                            b: sorted[1].choice_text,
                            c: sorted[2].choice_text,
                            d: sorted[3].choice_text
                        },
                        choiceIds: {
                            a: sorted[0].id,
                            b: sorted[1].id,
                            c: sorted[2].id,
                            d: sorted[3].id
                        }
                    };
                });
            } catch (err) {
                console.error('Fetch questions failed', err);
            }
        },
        async submitAnswers() {
            const payload = {
                user: this.user,
                answers: Object.entries(this.userAnswers).map(([question_id, choice_id]) => ({
                    question_id: parseInt(question_id),
                    choice_id
                }))
            };
            await fetch(`${this.apiHost}/answers`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });
        }
    }
};
</script>

<style scoped>
/* ===================== STYLING UNTUK Questions.vue ===================== */

/* .question-text: styling judul soal */
.question-text {
    font-size: 1.1rem;
    font-weight: 600;
    color: #ffffff;
    margin-bottom: 1rem;
    background: linear-gradient(135deg, #8b5cf6 0%, #d8b4fe 100%);
}

/* ===== PROGRESS BAR ===== */
.time-bar-container {
    height: 8px;
    background-color: #e5e7eb;
    border-radius: 4px;
    overflow: hidden;
}

.time-bar-fill {
    background-color: #fbbf24;
    height: 100%;
    transition: width 0.3s ease;
}

/* ===== LIST JAWABAN ===== */
.answer-list {
    display: flex;
    flex-direction: column;
}

.answer-option {
    background-color: #f3f4f6;
    border-radius: 0.75rem;
    padding: 0.75rem 1rem;
    margin-bottom: 0.75rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border: 2px solid transparent;
    transition: background-color 0.2s, border-color 0.2s;
}

.answer-option:hover {
    background-color: #e5e7eb;
}

/* Pilihan yang dipilih */
.answer-option.selected {
    background-color: #cffafe;
    border-color: #0ea5e9;
}

/* Huruf opsi (A, B, C, D) */
.option-label {
    font-weight: 600;
    font-size: 1rem;
    color: #1f2937;
    width: 24px;
    text-align: center;
}

/* Teks opsi */
.option-text {
    font-size: 0.95rem;
    color: #374151;
    flex: 1;
}
</style>