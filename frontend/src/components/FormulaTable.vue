
<template>
  <v-data-table
    :headers="headers"
    :items="formulaDataList"
    hide-actions
    disable-filtering
    class="elevation-5"
  >
    <template v-slot:items="props">
      <tr style="cursor: pointer;" v-on:click="navigateToFormula(props.item.id)">
        <td class="text-capitalize">{{ props.item.name }}</td>
        <td>{{ props.item.number }}</td>
        <td class="text-capitalize">{{ props.item.base }}</td>
        <td>{{ props.item.colorants }}</td>
        <td class="text-capitalize">{{ props.item.created }}</td>
      </tr>
    </template>
  </v-data-table>
</template>

<script lang="ts">
import Vue from "vue";
import { Formula } from "../basecoat_message_pb";
import * as moment from "moment";

import BasecoatClientWrapper from "../basecoatClientWrapper";

let client: BasecoatClientWrapper;
client = new BasecoatClientWrapper();

interface modifiedFormula {
  id: string;
  name: string;
  number: string;
  base: string;
  colorants: number;
  created: string;
}

export default Vue.extend({
  data: function() {
    return {
      headers: [
        {
          text: "Name",
          align: "left",
          value: "name"
        },
        {
          text: "Number",
          value: "number",
          sortable: false
        },
        {
          text: "Base",
          value: "base"
        },
        {
          text: "Colorants",
          value: "colorants"
        },
        {
          text: "Created",
          value: "created"
        }
      ],
      formulaDataList: [] as modifiedFormula[]
    };
  },
  created() {
    this.formulaDataList = this.formulaDataToList();
    this.$store.subscribe((mutation, state) => {
      if (
        mutation.type === "updateFormulaData" ||
        mutation.type === "updateFormulaDataFilter"
      ) {
        this.formulaDataList = this.formulaDataToList();
      }
    });
  },
  methods: {
    navigateToFormula: function(formulaID: string) {
      this.$router.push("/formulas/" + formulaID);
    },
    // This makes it so that the formula table is sortable.
    // the formula table sorts based on the data structure that
    // you pass it. So you have to pass it a data structure with
    // correct types in order of it to sort properly
    formulaDataToList(): modifiedFormula[] {
      let filteredIDs: string[] = this.$store.state.formulaDataFilter;

      let formulaDataMap: { [key: string]: Formula } = this.$store.state
        .formulaData;
      let formulaDataList: Formula[] = [];

      for (const [key, value] of Object.entries(formulaDataMap)) {
        if (filteredIDs && filteredIDs.length) {
          if (filteredIDs.includes(key)) {
            formulaDataList.push(value);
          }
        } else {
          formulaDataList.push(value);
        }
      }

      let modifiedFormulaList: modifiedFormula[] = [];
      let formula: Formula;

      for (formula of formulaDataList) {
        let modifiedFormula: modifiedFormula = {
          id: "",
          name: "",
          number: "",
          base: "",
          colorants: 0,
          created: ""
        };

        modifiedFormula.id = formula.getId();
        modifiedFormula.colorants = formula.getColorantsList().length;
        modifiedFormula.name = formula.getName();
        modifiedFormula.number = formula.getNumber();
        modifiedFormula.created = moment(
          moment.unix(formula.getCreated())
        ).format("L");
        modifiedFormula.base = "None";

        if (formula.getBasesList().length != 0) {
          modifiedFormula.base = formula.getBasesList()[0].getName();
        }
        modifiedFormulaList.push(modifiedFormula);
      }

      return modifiedFormulaList;
    }
  }
});
</script>

<style scoped>
</style>
