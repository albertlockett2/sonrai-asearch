import { randomUUID } from 'crypto';
import { setLogger } from 'grpc';
import { Edge, Filter, Search, SearchStep } from './gen/search_pb';

type LegacySearchFieldArgument = {
  caseSensitive?: boolean;
  id?: string;
  op: string;
  value: string | number | boolean | null;
};

type LegacySearchEdgeArgument = {
  count: LegacySearchFieldArgument;
  items: LegacySearchArguments;
};

interface LegacySearchArguments
  extends Record<
    string,
    | LegacySearchArguments[]
    | LegacySearchEdgeArgument
    | LegacySearchFieldArgument
  > {
  and: LegacySearchArguments[];
  or: LegacySearchArguments[];
}

type LegacySearchField = {
  childOperator: string;
  parent?: string;
  arguments: LegacySearchArguments;
  selectionSet: string[];
};

type LegacySearch = {
  fields: Record<string, LegacySearchField>;
};

type Context = {
  source: LegacySearch;
};

export function transformLegacy(source: LegacySearch): Search {
  const context: Context = { source };
  const search = new Search();
  const parentField = findParentField(source);
  if (parentField == null) {
    throw new Error('could not find parent field');
  }

  const step = toSearchStep(parentField, context);
  search.addSteps(step);

  return search;
}

function findParentField(source: LegacySearch): LegacySearchField | null {
  return (
    Object.values(source.fields).find((field) => {
      return !field.parent;
    }) ?? null
  );
}

/**
 * transform the field into a SearchStep
 * @param field
 * @param context
 * @returns {SearchStep}
 */
function toSearchStep(field: LegacySearchField, context: Context): SearchStep {
  const step = new SearchStep();
  step.setId(newStepId());
  step.setType(SearchStep.Type.FILTER);
  addFilters(step, field.arguments);
  addNextSteps(step, field, context);
  return step;
}

/**
 * adds filters to the step
 * @param step
 * @param argument
 */
function addFilters(step: SearchStep, argument: LegacySearchArguments): void {
  Object.keys(argument).forEach((property) => {
    const value = argument[property];

    // handle boolean logic elsewhere
    if ('and' === property || 'or' === property) {
      return;
    }

    // check to make sure it's not an edge filter
    if ((<any>value).items || (<any>value).count) {
      return;
    }

    const fieldArg = <LegacySearchFieldArgument>value;
    const filter = new Filter();
    filter.setProperty(property);
    // TODO handle other funny cases where it's not a string
    filter.setProperty(<string>fieldArg.value);
    step.addFilters(filter);
  });

  if (argument.and) {
    argument.and.forEach((and) => {
      addFilters(step, and);
    });
  }

  // TODO handle 'or' and 'not' fields
  if (argument.or) {
    console.warn('"or" arguments currently unsuported!!');
  }
}

/**
 * adds next steps to the step
 * @param step
 * @param field
 * @param context
 */
function addNextSteps(
  step: SearchStep,
  field: LegacySearchField,
  context: Context
): void {
  // TODO handle the case where the edge filters should be 'ANDed' not 'ORed'
  if (field.childOperator != 'or') {
    console.warn(
      `child operator "${field.childOperator}" not currently supported`
    );
    return;
  }

  const edges: Edge[] = [];
  Object.keys(field.arguments).forEach((edgeName) => {
    const value = field.arguments[edgeName];

    if ('and' === edgeName || 'or' === edgeName) {
      return;
    }

    // ensure we have correct kind of value
    const isEdgeField = (<any>value).items || (<any>value).count;
    if (!isEdgeField) {
      return;
    }

    // const edgeArg = <LegacySearchEdgeArgument>value;
    const edge = new Edge();
    // TODO really check the direction
    edge.setDir(Edge.Dir.OUT);
    edge.setName(edgeName);
    edges.push(edge);
  });

  // TODO handle case that there is more than one edge
  if (edges.length > 0) {
    const edgeStep = new SearchStep();
    edgeStep.setType(SearchStep.Type.EDGE);
    edgeStep.addEdges(edges[0]);

    const otherSideField = context.source.fields[field.selectionSet[0]];
    edgeStep.addNextsteps(toSearchStep(otherSideField, context));
    step.addNextsteps(edgeStep);
  }
}

/**
 * generate the ID of one of our steps
 */
function newStepId(): string {
  return randomUUID();
}
