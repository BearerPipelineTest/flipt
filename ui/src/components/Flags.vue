<template>
  <section class="section">
    <div class="container">
      <div class="level">
        <div class="level-left">
          <div class="level-item">
            <div v-if="flags.length > 10" class="field has-addons">
              <p class="control">
                <input
                  v-model="search"
                  class="input"
                  type="text"
                  placeholder="Find a Flag"
                />
              </p>
              <p class="control"><button class="button">Search</button></p>
            </div>
          </div>
        </div>
        <div class="level-right">
          <div class="level-item">
            <RouterLink
              data-testid="new-flag"
              class="button is-primary"
              :to="{ name: 'new-flag' }"
              >New Flag</RouterLink
            >
          </div>
        </div>
      </div>
      <b-table
        :data="isEmpty ? [] : filteredFlags"
        icon-pack="fas"
        :hoverable="true"
      >
        <b-table-column v-slot="props" field="enabled" label="Enabled">
          <span v-if="props.row.enabled" class="tag is-primary is-rounded"
            >On</span
          >
          <span v-else class="tag is-light is-rounded">Off</span>
        </b-table-column>
        <b-table-column v-slot="props" field="key" label="Key" sortable>
          <RouterLink :to="{ name: 'flag', params: { key: props.row.key } }">
            {{ props.row.key }}
          </RouterLink>
        </b-table-column>
        <b-table-column v-slot="props" field="name" label="Name" sortable>
          <RouterLink :to="{ name: 'flag', params: { key: props.row.key } }">
            {{ props.row.name }}
          </RouterLink>
        </b-table-column>
        <b-table-column v-slot="props" field="hasVariants" label="Variants">
          {{ props.row.variants.length > 0 ? "yes" : "no" }}
        </b-table-column>
        <b-table-column v-slot="props" field="description" label="Description">
          <small>{{ props.row.description | limit }}</small>
        </b-table-column>
        <b-table-column
          v-slot="props"
          field="createdAt"
          label="Created"
          sortable
        >
          <small>{{ props.row.createdAt | moment("from", "now") }}</small>
        </b-table-column>
        <b-table-column
          v-slot="props"
          field="updatedAt"
          label="Updated"
          sortable
        >
          <small>{{ props.row.updatedAt | moment("from", "now") }}</small>
        </b-table-column>
        <template #empty>
          <section class="section">
            <div class="content has-text-grey has-text-centered">
              <p>
                No flags found. Create a
                <RouterLink :to="{ name: 'new-flag' }">New Flag</RouterLink>.
              </p>
            </div>
          </section>
        </template>
      </b-table>
      <b-pagination
        v-if="totalCount > 20"
        v-model="currentPage"
        :per-page="1"
        :total="totalCount"
        @change="getFlags"
      >
      </b-pagination>
    </div>
  </section>
</template>

<script>
import { Api } from "@/services/api";
import notify from "@/mixins/notify";

export default {
  name: "Flags",
  mixins: [notify],
  data() {
    return {
      isEmpty: true,
      search: "",
      currentPage: 1,
      flags: [],
      total: 0,
    };
  },
  computed: {
    filteredFlags() {
      return this.flags.filter((flag) => {
        return (
          flag.name && flag.name.toLowerCase().match(this.search.toLowerCase())
        );
      });
    },
  },
  mounted() {
    this.getFlags();
  },
  methods: {
    getFlags() {
      Api.get("/flags?offset=" + (this.currentPage - 1))
        .then((response) => {
          this.isEmpty = false;
          this.flags = response.data.flags ? response.data.flags : [];
          this.total = response.data.total;
        })
        .catch((error) => {
          this.notifyError("Error loading flags.");
          console.error(error);
        });
    },
  },
};
</script>
