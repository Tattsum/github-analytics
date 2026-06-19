/* eslint-disable */
import * as types from './graphql';
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 * Learn more about it here: https://the-guild.dev/graphql/codegen/plugins/presets/preset-client#reducing-bundle-size
 */
type Documents = {
    "\n  query MemberDetail($login: String!) {\n    member(login: $login) {\n      login\n      name\n      totalCommits\n      totalPRCreated\n      totalPRMerged\n      totalReviews\n      totalIssues\n      totalAdditions\n      totalDeletions\n      prToReviewRatio\n      firstActivityYear\n      peakActivityYear\n      peakActivityCommits\n      yearlyStats {\n        year\n        commitCount\n        prCreated\n        prMerged\n        reviewCount\n        issueCount\n        totalAdditions\n        totalDeletions\n      }\n      topRepositories {\n        repository\n        commitCount\n        prCount\n        reviewCount\n        issueCount\n        totalAdditions\n        totalDeletions\n        firstActivity\n        lastActivity\n      }\n      longTermRepositories {\n        repository\n        commitCount\n        prCount\n        reviewCount\n        issueCount\n        totalAdditions\n        totalDeletions\n        firstActivity\n        lastActivity\n      }\n      roleTransition {\n        year\n        description\n        prCreated\n        reviewCount\n        ratio\n      }\n    }\n  }\n": typeof types.MemberDetailDocument,
    "\n  query Repositories {\n    repositories {\n      nameWithOwner\n      contributorCount\n      total {\n        commits\n        prCreated\n        prMerged\n        issues\n        reviews\n        additions\n        deletions\n      }\n    }\n  }\n": typeof types.RepositoriesDocument,
    "\n  query Repository($nameWithOwner: String!) {\n    repository(nameWithOwner: $nameWithOwner) {\n      nameWithOwner\n      contributorCount\n      total {\n        commits\n        prCreated\n        prMerged\n        issues\n        reviews\n        additions\n        deletions\n      }\n      contributors {\n        login\n        commitCount\n        prCreated\n        reviewCount\n        additions\n        deletions\n      }\n    }\n  }\n": typeof types.RepositoryDocument,
    "\n  query TeamOverviewSummary {\n    teamSummary {\n      memberCount\n      repositoryCount\n      totalCommits\n      totalPRCreated\n      totalPRMerged\n      totalIssues\n      totalReviews\n      totalAdditions\n      totalDeletions\n    }\n  }\n": typeof types.TeamOverviewSummaryDocument,
    "\n  query TeamOverviewMembers {\n    members {\n      login\n      name\n      totalCommits\n      totalPRCreated\n      totalPRMerged\n      totalIssues\n      totalReviews\n      totalAdditions\n      totalDeletions\n      prToReviewRatio\n    }\n  }\n": typeof types.TeamOverviewMembersDocument,
};
const documents: Documents = {
    "\n  query MemberDetail($login: String!) {\n    member(login: $login) {\n      login\n      name\n      totalCommits\n      totalPRCreated\n      totalPRMerged\n      totalReviews\n      totalIssues\n      totalAdditions\n      totalDeletions\n      prToReviewRatio\n      firstActivityYear\n      peakActivityYear\n      peakActivityCommits\n      yearlyStats {\n        year\n        commitCount\n        prCreated\n        prMerged\n        reviewCount\n        issueCount\n        totalAdditions\n        totalDeletions\n      }\n      topRepositories {\n        repository\n        commitCount\n        prCount\n        reviewCount\n        issueCount\n        totalAdditions\n        totalDeletions\n        firstActivity\n        lastActivity\n      }\n      longTermRepositories {\n        repository\n        commitCount\n        prCount\n        reviewCount\n        issueCount\n        totalAdditions\n        totalDeletions\n        firstActivity\n        lastActivity\n      }\n      roleTransition {\n        year\n        description\n        prCreated\n        reviewCount\n        ratio\n      }\n    }\n  }\n": types.MemberDetailDocument,
    "\n  query Repositories {\n    repositories {\n      nameWithOwner\n      contributorCount\n      total {\n        commits\n        prCreated\n        prMerged\n        issues\n        reviews\n        additions\n        deletions\n      }\n    }\n  }\n": types.RepositoriesDocument,
    "\n  query Repository($nameWithOwner: String!) {\n    repository(nameWithOwner: $nameWithOwner) {\n      nameWithOwner\n      contributorCount\n      total {\n        commits\n        prCreated\n        prMerged\n        issues\n        reviews\n        additions\n        deletions\n      }\n      contributors {\n        login\n        commitCount\n        prCreated\n        reviewCount\n        additions\n        deletions\n      }\n    }\n  }\n": types.RepositoryDocument,
    "\n  query TeamOverviewSummary {\n    teamSummary {\n      memberCount\n      repositoryCount\n      totalCommits\n      totalPRCreated\n      totalPRMerged\n      totalIssues\n      totalReviews\n      totalAdditions\n      totalDeletions\n    }\n  }\n": types.TeamOverviewSummaryDocument,
    "\n  query TeamOverviewMembers {\n    members {\n      login\n      name\n      totalCommits\n      totalPRCreated\n      totalPRMerged\n      totalIssues\n      totalReviews\n      totalAdditions\n      totalDeletions\n      prToReviewRatio\n    }\n  }\n": types.TeamOverviewMembersDocument,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = graphql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function graphql(source: string): unknown;

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query MemberDetail($login: String!) {\n    member(login: $login) {\n      login\n      name\n      totalCommits\n      totalPRCreated\n      totalPRMerged\n      totalReviews\n      totalIssues\n      totalAdditions\n      totalDeletions\n      prToReviewRatio\n      firstActivityYear\n      peakActivityYear\n      peakActivityCommits\n      yearlyStats {\n        year\n        commitCount\n        prCreated\n        prMerged\n        reviewCount\n        issueCount\n        totalAdditions\n        totalDeletions\n      }\n      topRepositories {\n        repository\n        commitCount\n        prCount\n        reviewCount\n        issueCount\n        totalAdditions\n        totalDeletions\n        firstActivity\n        lastActivity\n      }\n      longTermRepositories {\n        repository\n        commitCount\n        prCount\n        reviewCount\n        issueCount\n        totalAdditions\n        totalDeletions\n        firstActivity\n        lastActivity\n      }\n      roleTransition {\n        year\n        description\n        prCreated\n        reviewCount\n        ratio\n      }\n    }\n  }\n"): (typeof documents)["\n  query MemberDetail($login: String!) {\n    member(login: $login) {\n      login\n      name\n      totalCommits\n      totalPRCreated\n      totalPRMerged\n      totalReviews\n      totalIssues\n      totalAdditions\n      totalDeletions\n      prToReviewRatio\n      firstActivityYear\n      peakActivityYear\n      peakActivityCommits\n      yearlyStats {\n        year\n        commitCount\n        prCreated\n        prMerged\n        reviewCount\n        issueCount\n        totalAdditions\n        totalDeletions\n      }\n      topRepositories {\n        repository\n        commitCount\n        prCount\n        reviewCount\n        issueCount\n        totalAdditions\n        totalDeletions\n        firstActivity\n        lastActivity\n      }\n      longTermRepositories {\n        repository\n        commitCount\n        prCount\n        reviewCount\n        issueCount\n        totalAdditions\n        totalDeletions\n        firstActivity\n        lastActivity\n      }\n      roleTransition {\n        year\n        description\n        prCreated\n        reviewCount\n        ratio\n      }\n    }\n  }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query Repositories {\n    repositories {\n      nameWithOwner\n      contributorCount\n      total {\n        commits\n        prCreated\n        prMerged\n        issues\n        reviews\n        additions\n        deletions\n      }\n    }\n  }\n"): (typeof documents)["\n  query Repositories {\n    repositories {\n      nameWithOwner\n      contributorCount\n      total {\n        commits\n        prCreated\n        prMerged\n        issues\n        reviews\n        additions\n        deletions\n      }\n    }\n  }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query Repository($nameWithOwner: String!) {\n    repository(nameWithOwner: $nameWithOwner) {\n      nameWithOwner\n      contributorCount\n      total {\n        commits\n        prCreated\n        prMerged\n        issues\n        reviews\n        additions\n        deletions\n      }\n      contributors {\n        login\n        commitCount\n        prCreated\n        reviewCount\n        additions\n        deletions\n      }\n    }\n  }\n"): (typeof documents)["\n  query Repository($nameWithOwner: String!) {\n    repository(nameWithOwner: $nameWithOwner) {\n      nameWithOwner\n      contributorCount\n      total {\n        commits\n        prCreated\n        prMerged\n        issues\n        reviews\n        additions\n        deletions\n      }\n      contributors {\n        login\n        commitCount\n        prCreated\n        reviewCount\n        additions\n        deletions\n      }\n    }\n  }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query TeamOverviewSummary {\n    teamSummary {\n      memberCount\n      repositoryCount\n      totalCommits\n      totalPRCreated\n      totalPRMerged\n      totalIssues\n      totalReviews\n      totalAdditions\n      totalDeletions\n    }\n  }\n"): (typeof documents)["\n  query TeamOverviewSummary {\n    teamSummary {\n      memberCount\n      repositoryCount\n      totalCommits\n      totalPRCreated\n      totalPRMerged\n      totalIssues\n      totalReviews\n      totalAdditions\n      totalDeletions\n    }\n  }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query TeamOverviewMembers {\n    members {\n      login\n      name\n      totalCommits\n      totalPRCreated\n      totalPRMerged\n      totalIssues\n      totalReviews\n      totalAdditions\n      totalDeletions\n      prToReviewRatio\n    }\n  }\n"): (typeof documents)["\n  query TeamOverviewMembers {\n    members {\n      login\n      name\n      totalCommits\n      totalPRCreated\n      totalPRMerged\n      totalIssues\n      totalReviews\n      totalAdditions\n      totalDeletions\n      prToReviewRatio\n    }\n  }\n"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;