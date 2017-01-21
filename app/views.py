import json

from flask import render_template, request, jsonify
from app.basecoat import db_utils

from app import app, db, models



@app.route('/')
def index():
    formula_table = db_utils.get_table('Formula')
    formula_list = [formula for formula in formula_table]
    return render_template('index.html',
                           formula_list=formula_list)


@app.route('/formula/<int:formula_id>')
def get_formula(formula_id):
    formula = db_utils.get_object_from_table('Formula', 'id', formula_id)[0]
    colorant_list = json.loads(formula.colorants)
    base_list = json.loads(formula.bases)

    return render_template('view_formula.html',
                           formula=formula,
                           colorant_list=colorant_list,
                           base_list=base_list)


@app.route('/formula/add', methods=['GET', 'POST'])
def add_formula():
    if request.method == 'POST':
        form_data = request.json
        colorants = json.dumps(form_data.pop('colorant_list', None))
        bases = json.dumps(form_data.pop('base_list', None))

        form_data = {key: value.strip() for key, value in form_data.items()}

        if "formula_id" in form_data.keys():
            db_utils.update_db("Formula", "id", form_data['formula_id'], **form_data)
            db_utils.update_db("Formula", "id", form_data['formula_id'], colorants=colorants, bases=bases)

        else:
            new_formula = models.Formula(formula_name=form_data['formula_name'].title(),
                                         formula_number=form_data['formula_number'],
                                         customer_name=form_data['customer_name'].title(),
                                         colorants=colorants,
                                         bases=bases,
                                         summary=form_data['summary'],
                                         notes=form_data['notes'])


            try:
                db.session.add(new_formula)
                db.session.commit()
            except:
                db.session.rollback()
                raise

        return jsonify({'success':True}), 200
    else:
        return render_template('add_formula.html')


@app.route('/formula/edit/<int:formula_id>')
def edit_formula(formula_id):
    formula = db_utils.get_object_from_table('Formula', 'id', formula_id)[0]
    colorant_list = json.loads(formula.colorants)
    base_list = json.loads(formula.bases)
    return render_template('edit_formula.html',
                           formula=formula,
                           colorant_list=colorant_list,
                           base_list=base_list)


@app.route('/formula/delete/<int:formula_id>', methods=['DELETE'])
def delete_formula(formula_id):
    db_utils.delete_from_db('Formula', 'id', formula_id)
    return jsonify({'success':True}), 200
